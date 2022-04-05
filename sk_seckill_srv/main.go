package main

import (
	"sk_seckill_srv/conf"
	"sk_seckill_srv/database"
	"sk_seckill_srv/handler/predis"
	"sk_seckill_srv/handler/rabbitmq"
	"sk_seckill_srv/log"
)

func init() {
	//格外注意init顺序
	conf.Init("conf/config.json")

	log.Init()

	database.ConnInit()
	database.ModelInit()

	predis.Init()
}

//微服务架构
//func main() {
//	// New Service
//	service := grpc.NewService(
//		micro.Name("yuyoung.srv.sk_seckill_srv"),
//		micro.Version("latest"),
//	)
//
//	// Initialise service
//	service.Init()
//
//	// Register Handler
//	yuyoung_srv_sk_seckill_srv.RegisterSeckillServiceHandler(service.Server(), new(handler.SeckillImpl))
//
//	// Run service
//	if err := service.Run(); err != nil {
//		panic(err)
//	}
//}

func main() {
	rabbitmq.GetSkMqHandler()
}
