package main

import (
	"encoding/json"
	"net/http"
	"telegram-assistant-bot/models"
	"telegram-assistant-bot/pkg/bot"
	"telegram-assistant-bot/pkg/setting"
)

const SourceIni = "config/app.ini"

func init() {
	//配置加载
	setting.Setup(SourceIni)
	//数据库加载
	models.SetUp()
	//redis加载
	//err := gredis.Setup()
	//if err != nil {
	//    log.Println(err)
	//}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		data["code"] = "200"
		jsonStr, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(jsonStr))
	})
	bot.SetUp()
}
