package user

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_user_srv "sk_user_srv/proto/user"
	. "sk_web/utils"
)

func LoginHandler(ctx *gin.Context) {
	email := ctx.PostForm("username")
	password := ctx.PostForm("password")

	//校验密码是否正确
	client := grpc.NewService()
	service := yuyoung_srv_sk_user_srv.NewUserService("yuyoung.srv.sk_user_srv", client.Client())
	resp, err := service.Login(ctx, &yuyoung_srv_sk_user_srv.GeneralRequest{
		User: &yuyoung_srv_sk_user_srv.BasicUserInfo{
			Email:    email,
			Password: password,
		},
	})
	if err != nil {
		logrus.WithField("module", "login_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}
	token, err := GenerateToken(email, UserExpireDuration, UserSecretKey)
	if err != nil {
		logrus.WithField("module", "login_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "生成token失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":     resp.Code,
		"msg":      resp.Msg,
		"username": email,
		"token":    token,
	})
	logrus.WithField("email", email).WithField("password", password).Info("登陆成功")
}

func AdminLoginHandler(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	client := grpc.NewService()
	service := yuyoung_srv_sk_user_srv.NewUserService("yuyoung.srv.sk_user_srv", client.Client())
	resp, err := service.AdminLogin(ctx, &yuyoung_srv_sk_user_srv.GeneralRequest{
		User: &yuyoung_srv_sk_user_srv.BasicUserInfo{
			Email:    username,
			Password: password,
		},
	})
	if err != nil {
		logrus.WithField("module", "login_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "密码或用户名错误",
		})
		return
	}

	token, err := GenerateToken(username, AdminExpireDuration, AdminSecretKey)
	if err != nil {
		logrus.WithField("module", "admin_login_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "生成token失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":        resp.Code,
		"msg":         resp.Msg,
		"username":    username,
		"admin_token": token,
	})
	logrus.WithField("username", username).WithField("password", password).Info("管理员登陆成功")
}
