package main

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"sk_user_srv/conf"
	"sk_user_srv/database"
	"sk_user_srv/handler"
	"sk_user_srv/log"
	yuyoung_srv_sk_user_srv "sk_user_srv/proto/user"
)

func init() {
	//格外注意init顺序
	conf.Init("conf/config.json")

	log.Init()

	database.ConnInit()
	database.ModelInit()
}

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("yuyoung.srv.sk_user_srv"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	yuyoung_srv_sk_user_srv.RegisterUserServiceHandler(service.Server(), new(handler.UserImpl))

	// Run service
	if err := service.Run(); err != nil {
		panic(err)
	}
}
