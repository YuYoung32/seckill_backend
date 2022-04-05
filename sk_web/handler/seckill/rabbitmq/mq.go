package rabbitmq

import (
	"github.com/sirupsen/logrus"
	"sk_web/conf"
	"streadway/streadway/amqp"
)

var Mq *RabbitMQ

func MQInit() {
	Mq = NewRabbitMQ(
		"yuyoung.web.sk_web-order_queue",
		"yuyoung.web.sk_web-order_exchange",
		"direct",
		"yuyoung.web.sk_web-order_routing_key").Init()
}

type RabbitMQ struct {
	Conn         *amqp.Connection
	Ch           *amqp.Channel
	QueueName    string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

func NewRabbitMQ(queueName, exchangeName, exchangeType, routingKey string) *RabbitMQ {
	return &RabbitMQ{
		QueueName:    queueName,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
		RoutingKey:   routingKey,
	}
}

func (r *RabbitMQ) connMQ() {
	mqConf := conf.GetRabbitmqConf()
	dialer := "amqp://" + mqConf.User + ":" + mqConf.Password + "@" + mqConf.Host + ":" + mqConf.Port + "/" + mqConf.Vhost
	conn, err := amqp.Dial(dialer)
	if err != nil {
		logrus.WithField("err", err).Errorf("连接rabbitmq失败")
		return
	}
	r.Conn = conn
	logrus.Info("连接rabbitmq成功")

}

func (r *RabbitMQ) closeConn() {
	err := r.Conn.Close()
	if err != nil {
		logrus.WithField("err", err).Error("关闭rabbitmq失败")
		return
	}
	logrus.Info("关闭rabbitmq成功")
}

func (r *RabbitMQ) openChan() {
	ch, err := r.Conn.Channel()
	if err != nil {
		logrus.WithField("err", err).Error("获取channel失败")
		return
	}
	r.Ch = ch
	logrus.Info("获取channel成功")
}

func (r *RabbitMQ) closeChan() {
	err := r.Ch.Close()
	if err != nil {
		logrus.WithField("err", err).Error("关闭channel失败")
		return
	}
	logrus.Info("关闭channel成功")
}

func (r *RabbitMQ) publishInit() {
	_, err := r.Ch.QueueDeclare(r.QueueName, true, false, false, false, nil)
	if err != nil {
		logrus.WithField("err", err).Error("声明队列失败")
		return
	}

	err = r.Ch.ExchangeDeclare(r.ExchangeName, r.ExchangeType, true, false, false, false, nil)
	if err != nil {
		logrus.WithField("err", err).Error("声明交换机失败")
		return
	}

	err = r.Ch.QueueBind(r.QueueName, r.RoutingKey, r.ExchangeName, false, nil)
	if err != nil {
		logrus.WithField("err", err).Error("绑定队列失败")
		return
	}
}

func (r *RabbitMQ) PublishMsg(msg string) {
	err := r.Ch.Publish(r.ExchangeName, r.RoutingKey, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(msg),
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		logrus.WithField("module", "publish_msg").Error("发送Mq消息失败:" + err.Error())
		return
	}
	logrus.Debug("发送消息成功")
}

func (r *RabbitMQ) Init() *RabbitMQ {
	r.connMQ()
	r.openChan()
	r.publishInit()
	return r
}

func (r *RabbitMQ) Close() {
	r.closeChan()
	r.closeConn()
}
