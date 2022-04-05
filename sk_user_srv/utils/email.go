package utils

import (
	"fmt"
	"github.com/astaxie/beego/utils"
	"math/rand"
	"strings"
	"time"
)

//GenEmailCode 生成指定位数随机验证码
func GenEmailCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {

		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

//SendEmail 发送邮件
func SendEmail(to_email string, msg string, username string, password string, host string, port string) error {
	emailConfig := fmt.Sprintf(`{"username":"%s","password":"%s","host":"%s","port":%s}`, username, password, host, port)
	emailConn := utils.NewEMail(emailConfig) // beego下的

	emailConn.From = strings.TrimSpace(username)
	emailConn.To = []string{strings.TrimSpace(to_email)}
	emailConn.Subject = "注册验证码"
	//注意这里我们发送给用户的是激活请求地址
	emailConn.Text = msg

	err := emailConn.Send()
	return err
}
