package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_user_srv "sk_user_srv/proto/user"
	"sk_web/utils"
)

func RegisterHandler(ctx *gin.Context) {
	email := ctx.PostForm("email")
	code := ctx.PostForm("captche")
	password := ctx.PostForm("password")
	repassword := ctx.PostForm("repassword")
	if !utils.VerifyEmailFormat(email) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusNotFound,
			"msg":  "邮箱格式不正确",
		})
		return
	}

	if len(password) < 6 || len(password) > 16 || password != repassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusNotFound,
			"msg":  "两次密码不一致或密码格式不对",
		})
		return
	}

	client := grpc.NewService()
	service := yuyoung_srv_sk_user_srv.NewUserService("yuyoung.srv.sk_user_srv", client.Client())
	resp, err := service.Register(context.TODO(), &yuyoung_srv_sk_user_srv.RegisterUserRequest{
		User: &yuyoung_srv_sk_user_srv.UserInfo{
			BasicInfo: &yuyoung_srv_sk_user_srv.BasicUserInfo{
				Email:    email,
				Password: password,
			},
			Username:    email, //这里设置为邮箱
			Description: "normal user",
			Status:      "1",
		},
		Code: code,
	})
	if err != nil {
		logrus.WithField("module", "register_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}
