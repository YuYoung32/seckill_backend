package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-web"
	"github.com/sirupsen/logrus"
	"sk_web/conf"
	"sk_web/handler/seckill/predis"
	"sk_web/handler/seckill/rabbitmq"
	"sk_web/log"
	"sk_web/middleware"
	"sk_web/router"
)

func init() {
	conf.Init("conf/config.json")
	log.Init()
	rabbitmq.MQInit()
	predis.Init()
}

func main() {

	engine := gin.New()
	engine.Use(middleware.LoggerToFile("log/run.log"))
	engine.Use(middleware.CrosMiddleware)
	service := web.NewService(
		web.Name("yuyoung.web.sk_web"),
		web.Version("latest"),
		web.Handler(engine),
		web.Address(":8081"),
	)

	if err := service.Init(); err != nil {
		logrus.Fatal(err)
	}

	router.InitRouter(engine)

	if err := service.Run(); err != nil {
		logrus.Fatal(err)
	}
}
