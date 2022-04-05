//创建管理员账号的程序，需要时创建
package database

import (
	"sk_user_srv/conf"
	"sk_user_srv/utils"
	"strconv"
	"testing"
	"time"
)

//添加一些用户信息
func Test_AddUserAccount(t *testing.T) {
	conf.Init("../conf/config.json")
	ConnInit()
	conn := GetDBConn()
	hashedPassword, _ := utils.Encryption("123456")

	for i := 0; i < 200; i++ {
		email := "testuser" + strconv.Itoa(i) + "_" + strconv.Itoa(time.Now().Nanosecond()) + "@qq.com"
		user := new(User)
		user.Username = email
		user.Email = email
		user.Password = hashedPassword
		user.Description = "test_user"
		user.Status = "1"
		conn.Create(user)
	}
}

func Test_AddAdminAccount(t *testing.T) {
	//在此添加单个admin账号
	username := "skadmin" //重要：Username必须以 skadmin 开头
	password := "123456"

	admin := new(Admin)
	admin.Username = username
	hashedPassword, err := utils.Encryption(password)
	if err != nil {
		t.Error(err)
	}
	admin.Password = hashedPassword

	conf.Init("../conf/config.json")
	ConnInit()
	conn := GetDBConn()
	if res := conn.First(admin, "username = ?", username); res.RowsAffected < 1 {
		conn.Create(admin)
	}
}
