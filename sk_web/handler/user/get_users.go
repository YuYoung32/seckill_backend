package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_user_srv "sk_user_srv/proto/user"
	"strconv"
)

func GetUsersHandler(ctx *gin.Context) {
	var err error
	qcurrentPage := ctx.DefaultQuery("currentPage", "1")
	currentPage, err := strconv.Atoi(qcurrentPage)
	qpageSize := ctx.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(qpageSize)
	if err != nil {
		logrus.WithField("module", "get_user_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}
	start := (currentPage - 1) * pageSize

	service := grpc.NewService()
	client := yuyoung_srv_sk_user_srv.NewUserService("yuyoung.srv.sk_user_srv", service.Client())
	resp, err := client.GetUserInfo(context.TODO(), &yuyoung_srv_sk_user_srv.GetUserInfoRequest{
		Start:  int32(start),
		Amount: int32(pageSize),
	})
	if err != nil {
		logrus.WithField("module", "get_user_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.General.Code,
			"msg":  resp.General.Msg,
		})
		return
	}

	type FrontUsers struct {
		Email      string `json:"email"`
		Desc       string `json:"desc"`
		Status     string `json:"status"`
		CreateTime string `json:"create_time"`
	}
	var frontUsers []FrontUsers
	for _, user := range resp.User {
		frontUsers = append(frontUsers, FrontUsers{
			Email:      user.BasicInfo.Email,
			Desc:       user.Description,
			Status:     user.Status,
			CreateTime: user.CreateTime,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":         resp.General.Code,
		"msg":          resp.General.Msg,
		"front_users":  frontUsers,
		"total":        resp.Total,
		"current_page": currentPage,
		"page_size":    pageSize,
	})
}
