package main

import (
	"log"
	"telegram-assistant-bot/models"
	"telegram-assistant-bot/pkg/bot"
	"telegram-assistant-bot/pkg/gredis"
	"telegram-assistant-bot/pkg/setting"
)

const SourceIni = "config/app.ini"

func init() {
	//配置加载
	setting.Setup(SourceIni)
	//数据库加载
	models.SetUp()
	//redis加载
	err := gredis.Setup()
	if err != nil {
	   log.Println(err)
	}
}

func main() {
	bot.SetUp()
}
