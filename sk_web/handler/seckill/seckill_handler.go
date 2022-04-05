package seckill

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sk_web/handler/seckill/rabbitmq"
	. "sk_web/utils"
)

func FrontSeckillHandler(ctx *gin.Context) {
	id := ctx.PostForm("id") //活动ID
	//在middleware.UserAuth中获取用户email
	userEmail, ok := ctx.Get("userEmail")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "请先登录",
		})
		return
	}

	////原始方法，直接传到后端进行处理
	//service := grpc.NewService()
	//client := yuyoung_srv_sk_seckill_srv.NewSeckillService("yuyoung.srv.sk_seckill_srv", service.Client())
	//resp, err := client.FrontSeckill(ctx, &yuyoung_srv_sk_seckill_srv.SeckillRequest{
	//	Id:    id, //活动ID
	//	Email: userEmail.(string),
	//})
	//if resp.Code != "200" {
	//	logrus.WithField("module", "front_seckill_handler").Error(err)
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": resp.Code,
	//		"msg":  resp.Msg,
	//	})
	//	return
	//}

	//新方法传到rabbitmq进行处理
	rabbitmq.Mq.PublishMsg(MapToStr(map[string]interface{}{
		"eventId": id,
		"email":   userEmail.(string),
	}))

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"msg":  "下单中...",
	})

}
