package main

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"sk_product_srv/conf"
	"sk_product_srv/database"
	"sk_product_srv/handler"
	"sk_product_srv/log"
	yuyoung_srv_sk_product_srv "sk_product_srv/proto"
	yuyoung_srv_sk_event_srv "sk_product_srv/proto/event"
	"sync"
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
	productService := grpc.NewService(
		micro.Name("yuyoung.srv.sk_product_srv"),
		micro.Version("latest"),
	)

	// Initialise service
	productService.Init()

	// Register Handler
	yuyoung_srv_sk_product_srv.RegisterProductServiceHandler(productService.Server(), new(handler.ProductImpl))

	eventService := grpc.NewService(micro.Name("yuyoung.srv.sk_event_srv"), micro.Version("latest"))
	eventService.Init()
	err := yuyoung_srv_sk_event_srv.RegisterEventServiceHandler(eventService.Server(), new(handler.EventImpl))
	if err != nil {
		panic(err)
	}

	wait := sync.WaitGroup{}
	wait.Add(1)
	// Run service
	go func() {
		if err := productService.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := eventService.Run(); err != nil {
			panic(err)
		}
	}()
	wait.Wait()

}
