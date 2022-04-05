package rabbitmq

import (
	"sk_web/conf"
	"sk_web/utils"
	"testing"
)

func TestRabbitMQ_PublishMsg(t *testing.T) {
	conf.Init("../conf/config.json")

	mq := NewRabbitMQ("testQueue", "testExchange", "direct", "testRoutingKey").Init()

	mq.PublishMsg(utils.MapToStr(map[string]interface{}{
		"name": "test",
		"age":  18,
	}))

	mq.Close()
}
