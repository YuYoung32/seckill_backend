package router

import (
	"github.com/gin-gonic/gin"
	"sk_web/handler/seckill"
	"sk_web/middleware"
)

func SeckillRouter(router *gin.RouterGroup) {
	router.POST("/front/seckill", middleware.UserAuth, seckill.FrontSeckillHandler)
	router.GET("/front/seckill_result", middleware.UserAuth, seckill.FrontSeckillResultHandler)
}
