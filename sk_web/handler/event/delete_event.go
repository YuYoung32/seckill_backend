package event

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_event_srv "sk_product_srv/proto/event"
)

func DeleteEventHandler(ctx *gin.Context) {
	id := ctx.PostForm("id")

	service := grpc.NewService()
	client := yuyoung_srv_sk_event_srv.NewEventService("yuyoung.srv.sk_event_srv", service.Client())
	resp, err := client.DeleteEvent(ctx, &yuyoung_srv_sk_event_srv.GeneralRequest{
		EventId: id,
	})
	if err != nil {
		logrus.WithField("module", "delete_event_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "活动不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}
