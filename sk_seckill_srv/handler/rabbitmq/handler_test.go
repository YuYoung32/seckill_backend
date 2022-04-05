package rabbitmq

import (
	"github.com/sirupsen/logrus"
	"sk_seckill_srv/conf"
	"testing"
)

func TestGetSkMqHandler(t *testing.T) {
	conf.Init("../../conf/config.json")
	r := new(RabbitMQ).Init()

	deliveries, err := r.Ch.Consume(
		"yuyoung.web.sk_web-order_queue", //队列名称
		"yuyoung.srv.sk_seckill_srv-order",
		false,
		false,
		false,
		false,
		nil,
	)
	//r.Ch.Qos(10, 0, false)
	if err != nil {
		logrus.WithField("module", "sk_mq_handler").Error(err)
		return
	}

	//从chan中读取消息
	for {
		select {
		case delivery := <-deliveries:
			go SkMqHandle(delivery)
		}
	}

}
