package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	userGroup := router.Group("/user")
	UserRouter(userGroup)

	productGroup := router.Group("/product")
	ProductRouter(productGroup)

	eventGroup := router.Group("/event")
	EventRouter(eventGroup)

	seckillGroup := router.Group("/seckill")
	SeckillRouter(seckillGroup)
}
