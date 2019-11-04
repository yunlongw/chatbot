package models_test

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"reflect"
	"telegram-assistant-bot/models"
	"telegram-assistant-bot/pkg/setting"
	"testing"
	"time"
)

const SourceIni = "../config/app.ini"
const ChatId  = -5656565

func init() {
	//配置加载
	setting.Setup(SourceIni)
	//数据库加载
	models.SetUp()
}

func TestAddGroup(t *testing.T) {
	var groupId int64 = -22222
	group := make(map[string]interface{})
	group["group_id"] = groupId
	group["title"] = "qweqwe"
	err := models.AddGroup(group)
	if err != nil {
		t.Error(err)
	}
}

func TestGetGroups(t *testing.T) {
	maps := make(map[string]interface{})
	_, err := models.GetGroups(0, 10, maps)
	if err != nil {
		t.Error(err)
	}
}

func TestExistGroups(t *testing.T) {
	_, err := models.ExistGroups(1)
	if err != nil {
		t.Error(err)
	}
}

func TestExistGroupsByGroupId(t *testing.T) {
	_, err := models.ExistGroups(11)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTotalGroup(t *testing.T) {
	maps := make(map[string]interface{})
	_, err := models.GetTotalGroup(maps)
	if err != nil {
		t.Error(err)
	}
}

func TestTotalCheck(t *testing.T) {
	maps := make(map[string]interface{})
	groups, err := models.GetGroups(0, 1000, maps)
	if err != nil {
		t.Error(err)
	}
	var total int
	total, err = models.GetTotalGroup(maps)
	if err != nil {
		t.Error(err)
	}

	groupsLen := len(groups)

	if groupsLen > total {
		t.Error("数量不一致")
	}
}

func TestSetGroupSetting(t *testing.T) {
	b, err := models.SetGroupSetting(ChatId, fmt.Sprintf("%d", time.Now().YearDay()), fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		t.Error(err)
	}

	if b == false {
		t.Error("添加失败")
	}
}

func TestGetGroupSettings(t *testing.T) {
	maps := make(map[string]interface{})
	maps["group_id"] = ChatId
	_, err := models.GetGroupSettings(maps)
	if err != nil {
		t.Error(err)
	}
}

func TestGetGroupSettingByGroupIDToHashMap(t *testing.T) {
	_ , err := models.GetGroupSettingByGroupIDToHashMap(ChatId)
	if err != nil {
	    t.Error(err)
	}
}

func TestSetUp(t *testing.T) {
	var number3 = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
		tgbotapi.NewInlineKeyboardButtonData("test", "test"),
	)

	fmt.Println("type:", reflect.TypeOf(number3))
}
