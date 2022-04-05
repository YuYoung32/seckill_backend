package event

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_event_srv "sk_product_srv/proto/event"
	"strconv"
)

func EditEventHandler(ctx *gin.Context) {
	id := ctx.PostForm("id")
	name := ctx.PostForm("name")
	pprice := ctx.PostForm("price")
	pnum := ctx.PostForm("num")
	pid := ctx.PostForm("pid")
	startTime := ctx.PostForm("start_time")
	endTime := ctx.PostForm("end_time")
	if id == "" || name == "" || pprice == "" || pnum == "" || pid == "" || startTime == "" || endTime == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数不能为空",
		})
		return
	}

	price, err := strconv.ParseFloat(pprice, 32)
	if err != nil {
		logrus.WithField("module", "edit_event_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "价格格式错误",
		})
		return
	}
	num, err := strconv.Atoi(pnum)
	if err != nil {
		logrus.WithField("module", "edit_event_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "数量格式错误",
		})
		return
	}

	service := grpc.NewService()
	client := yuyoung_srv_sk_event_srv.NewEventService("yuyoung.srv.sk_event_srv", service.Client())
	resp, err := client.EditEvent(context.TODO(), &yuyoung_srv_sk_event_srv.EditEventRequest{
		EventInfo: &yuyoung_srv_sk_event_srv.EventInfo{
			Id:         id,
			Name:       name,
			EventPrice: float32(price),
			EventNum:   int32(num),
			StartTime:  startTime,
			EndTime:    endTime,
			ProductId:  pid,
		},
	})
	if err != nil {
		logrus.WithField("module", "edit_event_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}
