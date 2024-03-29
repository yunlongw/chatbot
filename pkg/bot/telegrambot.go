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
	"telegram-assistant-bot/pkg/gredis"
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
	setting_Verify = "DefaultVerify"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		data["code"] = "200"
		id := captcha.NewLen(4)
		url := setting.BotSetting.HttpServer + "captcha/" + id + ".png"
		data["url"] = url
		jsonStr, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(jsonStr))
	})
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))

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
			callbackQuery(update)
		}

		if update.Message == nil {
			log.Println("error...")
			continue
		}

		// 判断常规信息
		if update.Message != nil {
			//验证信息
			verifyAction(update)
			//captcha.VerifyString()
			log.Printf("%s", update.Message.Text)
		}

		// 检测加入分组和离开分组
		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			if update.Message.NewChatMembers != nil {
				newChatMembers(update)
			}

			// 离开分组判断
			if update.Message.LeftChatMember != nil {
				leftChatMember(update)
			}
		}

		// 判断 command
		if update.Message.IsCommand() {
			doCommand(update)
		}

	}
}

func verifyAction(update tgbotapi.Update) {
	if id, err := getVerify(*update.Message.From); err == nil {
		code := update.Message.Text
		log.Printf("id:%s", id)
		log.Printf("code:%s", []byte(code))
		if captcha.VerifyString(id, code) {
			fmt.Println("成功")
			ch.Chan(update.Message.Chat.ID,update.Message.From.ID, true)
		} else {
			fmt.Println("失败")
			ch.Chan(update.Message.Chat.ID,update.Message.From.ID, false)
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
		create, administrators := adminList(update.Message.Chat.ID)
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

func adminList(ChatID int64) (tgbotapi.User, []*tgbotapi.User) {
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

func leftChatMember(update tgbotapi.Update) (tgbotapi.Message, error) {
	return bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s left this group, Bye,Bye!", update.Message.LeftChatMember.UserName)))
}

const MessageModel = `通过复选按钮，调整设置。提醒：建议看官网首页对相关功能的更详细说明。

刚刚进行的更改：%s

推荐设置：启用审核并信任管理，不使用记录模式。静音模式避免打扰其他人，私信设置让机器人通过私聊发送设置菜单。`

func callbackQuery(update tgbotapi.Update) {
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
		ok, err := models.SetGroupSetting(update.CallbackQuery.Message.Chat.ID, setting_Verify, "1")
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
		ok, err := models.SetGroupSetting(update.CallbackQuery.Message.Chat.ID, setting_Verify, "0")
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
func newChatMembers(update tgbotapi.Update) {
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

					user, _ := adminList(update.Message.Chat.ID)
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

			if user.IsBot {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,"抓到一只机器人")
				_, _ = bot.Send(msg)
				return
			}

			verifyData(update, user)

		}
	}
}

// 验证数据
func verifyData(update tgbotapi.Update, user tgbotapi.User) {
	id := captcha.NewLen(6)
	url := setting.BotSetting.HttpServer + "captcha/" + id + ".png"
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: content}
	messageWithPhoto := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, bytes)
	messageWithPhoto.Caption = fmt.Sprintf(Verfily, "@"+getUserName(user), 6, 50)
	messageWithPhoto.ReplyMarkup = nKeyboard1
	msg, err := bot.Send(messageWithPhoto)
	if err != nil {
		log.Println(err)
	}
	setVerify(user, id)
	ch.SetChan(update.Message.Chat.ID,user.ID)
	go func() {
		select {
		case <-time.After(50 * time.Second):
			fmt.Println("超时测试")
			sendMessage(tgbotapi.NewEditMessageCaption(update.Message.Chat.ID, msg.MessageID, "超时验证 timeout"))
			deleteVerify(user,id)
			ch.DeleteChan(update.Message.Chat.ID, user.ID)
			//todo 移除对话框
			deleteVerifyMessage(update.Message.Chat.ID, update.Message.MessageID)
		case m := <- ch.m[update.Message.Chat.ID][user.ID]:
			if m == true {
				sendMessage(tgbotapi.NewEditMessageCaption(update.Message.Chat.ID, msg.MessageID, "验证通过"))
			}else {
				sendMessage(tgbotapi.NewEditMessageCaption(update.Message.Chat.ID, msg.MessageID, "验证失败"))
			}
			deleteVerify(user,id)
			ch.DeleteChan(update.Message.Chat.ID, user.ID)
			deleteVerifyMessage(update.Message.Chat.ID, update.Message.MessageID)
		}
	}()
}

func deleteVerifyMessage(chatID int64, messageID int)  {
	go func() {
		select {
			case <- time.After(5* time.Second):
				sendMessage(tgbotapi.NewDeleteMessage(chatID, messageID))
		}
	}()
}

const PassMessage =`您已通过验证。
单码用时：%d 秒
You have passed the CAPTCHA test.
Time used: %d seconds.`

var ch = NewGroupVerifyChanMap()

func setVerify(user tgbotapi.User, id string) bool {
	var key = fmt.Sprintf("verify:%d", user.ID)
	err := gredis.Set(key, id, 50)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func deleteVerify(user tgbotapi.User, id string) bool {
	var key = fmt.Sprintf("verify:%d", user.ID)
	b, err := gredis.Delete(key)
	if err != nil {
		log.Println(err)
		return false
	}
	return b
}


func getVerify(user tgbotapi.User) (string, error) {
	var v interface{}
	var key = fmt.Sprintf("verify:%d", user.ID)
	result, err := gredis.Get(key)
	json.Unmarshal(result, &v)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if str, ok := v.(string); ok == true {
		return str, nil
	}
	return "", nil
}

var nKeyboard1 = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("看不清，换一个", "reset"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("人工通过", "en"),
		tgbotapi.NewInlineKeyboardButtonData("人工拒绝", "di"),
	),
)

const Verfily = `
针对 %s 的验证码
您好，请在 75 秒内输入上图 %d 英文字符验证码(不限大小写，不含数字)。
Welcome, please input 5-char CAPTCHA above in %d seconds.
`

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
		return fmt.Sprintf("%s", user.FirstName)
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

	if maps[setting_Verify] == "0" {
		buttonrows[0] = append(buttonrows[0], tgbotapi.NewInlineKeyboardButtonData("□审核开关", "enable"))
	} else {
		buttonrows[0] = append(buttonrows[0], tgbotapi.NewInlineKeyboardButtonData("■审核开关", "disable"))
	}
	buttonrows[1] = append(buttonrows[1], tgbotapi.NewInlineKeyboardButtonData("测试1", "test11"))
	buttonrows[2] = append(buttonrows[2], tgbotapi.NewInlineKeyboardButtonData("<<back", "test"))

	return tgbotapi.NewInlineKeyboardMarkup(buttonrows...)
}
