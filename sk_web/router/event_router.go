package router

import (
	"github.com/gin-gonic/gin"
	"sk_web/handler/event"
	"sk_web/middleware"
)

func EventRouter(router *gin.RouterGroup) {
	router.GET("/get_products", middleware.AdminAuth, event.GetProductsHandler)
	router.GET("/get_events", middleware.AdminAuth, event.GetEventsHandler)
	router.GET("/get_event", middleware.AdminAuth, event.GetEventHandler)
	router.POST("/add_event", middleware.AdminAuth, event.AddEventHandler)
	router.POST("/edit_event", middleware.AdminAuth, event.EditEventHandler)
	router.POST("/delete_event", middleware.AdminAuth, event.DeleteEventHandler)

	router.GET("/front/get_events", event.GetFrontEventsHandler)
	router.GET("/front/get_event", middleware.UserAuth, event.GetFrontEventHandler)
}
