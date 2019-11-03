package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type Bot struct {
	ChatID int
	ApiToken string
}

var BotSetting = &Bot{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DataBaseSetting = &Database{}
var cfg *ini.File

func Setup(source string) {
	var err error
	cfg, err = ini.Load(source)

	if err != nil {
		log.Println(err)
	}

	mapTo("database", DataBaseSetting)
	mapTo("bot", BotSetting)
	mapTo("redis", RedisSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
