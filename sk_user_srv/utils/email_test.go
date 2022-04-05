package utils

import (
	"fmt"
	"sk_user_srv/conf"
	"testing"
	"time"
)

func TestGenEmailCode(t *testing.T) {
	for width := 1; width <= 5; width++ {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Nanosecond * 10)
			code := GenEmailCode(width)
			fmt.Println(code)
		}
	}
}

func TestSendEmail(t *testing.T) {
	conf.Init("../conf/config.json")
	config := conf.GetEmailConf()
	err := SendEmail("3400711168@qq.com", "6666", config.Username, config.Password, config.Host, config.Port)
	if err != nil {
		t.Error(err)
	}
}
