package router

import (
	"github.com/gin-gonic/gin"
	"sk_web/handler/user"
	"sk_web/middleware"
)

func UserRouter(router *gin.RouterGroup) {
	router.POST("/send_email", user.SendEmailHandler)
	router.POST("/register", user.RegisterHandler)
	router.POST("/login", user.LoginHandler)
	router.POST("/admin_login", user.AdminLoginHandler)
	router.GET("/get_users", middleware.AdminAuth, user.GetUsersHandler)
}
