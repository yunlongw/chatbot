package bot

import (
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/url"
	"strings"
	"telegram-assistant-bot/models"
	"telegram-assistant-bot/pkg/setting"
	"time"
)

const (
	welcome  = `Welcome to test chat! We can help in English 🇬🇧`
	chatInfo = `
create: 
%s
admin:
%s
`
)

const (
	Setting_Verify = "DefaultVerify"
)

const (
	// Default number of digits in captcha solution.
	DefaultLen = 6
	// The number of captchas created that triggers garbage collection used
	// by default store.
	CollectNum = 100
	// Expiration time of captchas used by default store.
	Expiration = 10 * time.Minute
	// Standard width and height of a captcha image.
	StdWidth  = 240
	StdHeight = 80
)

var bot *tgbotapi.BotAPI

var numericKeyboard1 = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
		tgbotapi.NewInlineKeyboardButtonData("test", "test"),
	),
)
var numericKeyboard3 = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("setting", "setting"),
		tgbotapi.NewInlineKeyboardButtonData("sign", "sign"),
		tgbotapi.NewInlineKeyboardButtonData("<<back", "mean"),
	),
)

var numericKeyboard_setting = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("启用", "enable"),
		tgbotapi.NewInlineKeyboardButtonData("禁用", "disable"),
		tgbotapi.NewInlineKeyboardButtonData("<<back", "test"),
	),
)

var disable bool

var singList map[int][]string

func SetUp() {
	singList = make(map[int][]string)
	var err error
	bot, err = tgbotapi.NewBotAPI(setting.BotSetting.ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	//pic := captcha.NewLen(6)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		data["code"] = "200"
		jsonStr, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(jsonStr))
	})

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(setting.BotSetting.HttpServer + bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:8001", nil)

	messageDispose(updates)
}

// 消息处理
func messageDispose(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		printJson(update)

		// 判断回调 ballBackQuery
		if update.CallbackQuery != nil {
			CallbackQuery(update)
		}

		if update.Message == nil {
			log.Println("error...")
			continue
		}

		// 判断常规信息
		if update.Message != nil {
			log.Printf("%s", update.Message.Text)
		}

		// 检测加入分组和离开分组
		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			if update.Message.NewChatMembers != nil {
				NewChatMembers(update)
			}

			// 离开分组判断
			if update.Message.LeftChatMember != nil {
				LeftChatMember(update)
			}
		}

		// 判断 command
		if update.Message.IsCommand() {
			doCommand(update)
		}

	}
}

func doCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "mean":
		message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(MessageModel, "主菜单"))
		message.ReplyToMessageID = update.Message.From.ID
		message.ReplyMarkup = numericKeyboard1
		sendMessage(message)
	case "admin":
		create, administrators := AdminList(update.Message.Chat.ID)
		var users []string
		for _, v := range administrators {
			users = append(users, getUserName(*v))
		}
		userListString := strings.Join(users, " \n")
		str := fmt.Sprintf(chatInfo, getUserName(create), userListString)
		log.Println(str)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
		sendMessage(msg)
	}
}

func AdminList(ChatID int64) (tgbotapi.User, []*tgbotapi.User) {
	chatMember, err := bot.GetChatAdministrators(tgbotapi.ChatConfig{
		ChatID: ChatID,
	})
	if err != nil {
		log.Println(err)
	}
	var create tgbotapi.User
	var administrators []*tgbotapi.User
	for _, val := range chatMember {
		if val.Status == "creator" {
			create = *val.User
		}
		if val.Status == "administrator" {
			administrators = append(administrators, val.User)
		}
	}
	return create, administrators
}

func LeftChatMember(update tgbotapi.Update) (tgbotapi.Message, error) {
	return bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s left this group, Bye,Bye!", update.Message.LeftChatMember.UserName)))
}

const MessageModel = `通过复选按钮，调整设置。提醒：建议看官网首页对相关功能的更详细说明。

刚刚进行的更改：%s

推荐设置：启用审核并信任管理，不使用记录模式。静音模式避免打扰其他人，私信设置让机器人通过私聊发送设置菜单。`

func CallbackQuery(update tgbotapi.Update) {
	switch update.CallbackQuery.Data {
	case "mean":
		msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, numericKeyboard1)
		sendMessage(msg)
	case "test":
		msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, numericKeyboard3)
		sendMessage(msg)
	case "setting":
		set := getSettingNewInlineKeyboardMarkup(update)
		msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, set)
		sendMessage(msg)
	case "sign":
		_, signMsg := userSign(*update.CallbackQuery.From)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, signMsg)
		sendMessage(msg)
	case "enable":
		ok, err := models.SetGroupSetting(update.CallbackQuery.Message.Chat.ID, Setting_Verify, "1")
		if err != nil {
			fmt.Println(err)
		}
		if ok == true {
			set := getSettingNewInlineKeyboardMarkup(update)
			message := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf(MessageModel, "启用"))
			message.ReplyMarkup = &set
			sendMessage(message)
		}

	case "disable":
		ok, err := models.SetGroupSetting(update.CallbackQuery.Message.Chat.ID, Setting_Verify, "0")
		if err != nil {
			fmt.Println(err)
		}
		if ok == true {
			set := getSettingNewInlineKeyboardMarkup(update)
			message := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf(MessageModel, "禁用"))
			message.ReplyMarkup = &set
			sendMessage(message)
		}
	default:
		apiResponse, _ := bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
		printJson(apiResponse)
		sendMessage(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
	}
}

//  新用户处理
func NewChatMembers(update tgbotapi.Update) {
	//var newUsers []string
	for _, user := range *update.Message.NewChatMembers {
		// 加入的机器人本身
		if user.ID == setting.BotSetting.ChatID {
			if update.Message.Chat.Type == "group" {
				group := make(map[string]interface{})
				group["group_id"] = update.Message.Chat.ID
				group["title"] = update.Message.Chat.Title

				if ok, _ := models.ExistGroupsByGroupId(update.Message.Chat.ID); ok != true {
					err := models.AddGroup(group)
					if err != nil {
						log.Println(err)
					}

					user, _ := AdminList(update.Message.Chat.ID)
					log.Println(user)
				}

				maps := make(map[string]interface{})
				maps["group_id"] = update.Message.Chat.ID
				maps["admin_id"] = user.ID

				ok, err := models.ExistAdminsGroups(maps)
				if err != nil {
					log.Println(err)
				}

				if ok != true {
					err = models.AddAdminsGroups(maps)
					if err != nil {
						log.Println(err)
					}
				}

			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("我是机器人 %s, 很高兴为您服务!", getUserName(user)))
			_, _ = bot.Send(msg)
		} else {
			//TODO 发送验证码
			//TODO 超时设置
			//TODO 失败拒绝
			//TODO 成功通过
			//id := captcha.NewLen(4)
			//digits := captcha.RandomDigits(4)
			//captcha.NewImage( id, digits, 30, 30)
			//file, err := url.Parse("http://cdn2.jianshu.io/assets/default_avatar/12-aeeea4bedf10f2a12c0d50d626951489.jpg")
			//if err != nil {
			//    panic(err)
			//}

			id := captcha.NewLen(4)
			url := setting.BotSetting.HttpServer + "captcha/" + id + ".png"
			fmt.Println(url)
			//res, err := http.Get("http://cdn2.jianshu.io/assets/default_avatar/12-aeeea4bedf10f2a12c0d50d626951489.jpg")
			res, err := http.Get(url)

			if err != nil {
				panic(err)
			}

			content, err := ioutil.ReadAll(res.Body)
			if err != nil {
				// error handling...
			}
			bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: content}
			messageWithPhoto := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, bytes)
			sendMessage(messageWithPhoto)


			//newUsers = append(newUsers, "@"+getUserName(user))
			//joinedUsers := strings.Join(newUsers, " ")
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hey, %s\n%s", joinedUsers, welcome))
			//msg.ReplyMarkup = Keyboard
			//_, _ = bot.Send(msg)
		}
	}
}

var Keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("重新生成二维码", "recreate"),
	),
)

func sendMessage(msg tgbotapi.Chattable) bool {
	if msg, err := bot.Send(msg); err != nil {
		printJson(msg)
		log.Println(err)
		return false
	} else {
		printJson(msg)
		return true
	}
}

func userSign(user tgbotapi.User) (b bool, signMsg string) {
	today := time.Now().Format("20060102")
	formId := user.ID
	if isSet(singList[formId], today) {
		b = false
		signMsg = "重复签到"
		return
	}
	singList[formId] = append(singList[formId], today)
	count := len(singList[user.ID])
	UserName := "@" + getUserName(user)
	signMsg = fmt.Sprintf("%s 签到成功!积分: %d", UserName, count)
	return
}

func getUserName(user tgbotapi.User) string {
	if user.UserName == "" {
		return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	return user.UserName
}

func printJson(v interface{}) (s string) {
	str, _ := json.Marshal(v)
	s = string(str)
	log.Println(s)
	return
}

func isSet(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}

func getSettingNewInlineKeyboardMarkup(update tgbotapi.Update) tgbotapi.InlineKeyboardMarkup {
	maps, err := models.GetGroupSettingByGroupIDToHashMap(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		log.Println(err)
	}
	buttonrows := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttonrows = append(buttonrows, make([]tgbotapi.InlineKeyboardButton, 0))
	buttonrows = append(buttonrows, make([]tgbotapi.InlineKeyboardButton, 0))
	buttonrows = append(buttonrows, make([]tgbotapi.InlineKeyboardButton, 0))

	if maps[Setting_Verify] == "0" {
		buttonrows[0] = append(buttonrows[0], tgbotapi.NewInlineKeyboardButtonData("□审核开关", "enable"))
	} else {
		buttonrows[0] = append(buttonrows[0], tgbotapi.NewInlineKeyboardButtonData("■审核开关", "disable"))
	}
	buttonrows[1] = append(buttonrows[1], tgbotapi.NewInlineKeyboardButtonData("测试1", "test11"))
	buttonrows[2] = append(buttonrows[2], tgbotapi.NewInlineKeyboardButtonData("<<back", "test"))

	return tgbotapi.NewInlineKeyboardMarkup(buttonrows...)
}
