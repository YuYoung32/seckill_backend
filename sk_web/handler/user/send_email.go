package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_user_srv "sk_user_srv/proto/user"
	_ "sk_web/log"
	. "sk_web/utils"
)

func SendEmailHandler(ctx *gin.Context) {
	email := ctx.PostForm("email")
	//邮件地址格式验证
	if !VerifyEmailFormat(email) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusNotFound,
			"msg":  "邮箱格式不正确",
		})
		return
	}

	service := grpc.NewService()
	service.Init()
	client := yuyoung_srv_sk_user_srv.NewUserService("yuyoung.srv.sk_user_srv", service.Client())
	resp, err := client.SendEmail(context.TODO(), &yuyoung_srv_sk_user_srv.GeneralRequest{
		User: &yuyoung_srv_sk_user_srv.BasicUserInfo{
			Email: email,
		},
	})
	if err != nil {
		logrus.WithField("module", "send_email").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}

	logrus.WithField("email", email).Info("发送邮件成功")
	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}
