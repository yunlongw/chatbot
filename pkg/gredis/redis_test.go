package gredis_test

import (
	"encoding/json"
	"log"
	"telegram-assistant-bot/pkg/gredis"
	"telegram-assistant-bot/pkg/setting"
	"testing"
)

const SourceIni = "../../config/app.ini"
const Key = "test"
const Result = "ss"

func init() {
	//配置加载
	setting.Setup(SourceIni)
	//redis加载
	err := gredis.Setup()
	if err != nil {
		log.Println(err)
	}
}

func TestSet(t *testing.T) {
	str := Result
	err := gredis.Set(Key, str, 30)
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	result, err := gredis.Get(Key)
	if err != nil {
		t.Error(err)
	}
	var v interface{}
	json.Unmarshal(result, &v)
	v , ok:= v.(string)
	if ok!=true {
		t.Error("转换失败")
	}
	if v != Result {
		log.Println(v)
		t.Error("测试结果和预期不符")
	}

}

