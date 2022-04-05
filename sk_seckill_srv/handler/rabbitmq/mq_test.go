package rabbitmq

import (
	"fmt"
	"sk_seckill_srv/conf"
	"testing"
)

func TestRabbitMQ_PublishMsg(t *testing.T) {
	conf.Init("../conf/config.json")

	mq := new(RabbitMQ).Init()
	//声明一个队列
	delivery, err := mq.Ch.Consume(
		"first_queue", //队列名称
		"myConsumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("消费失败")
		return
	}
	for msg := range delivery {
		fmt.Println(string(msg.Body))
		msg.Ack(true)
	}

	mq.Close()
}
