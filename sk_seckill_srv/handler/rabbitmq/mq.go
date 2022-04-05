package rabbitmq

import (
	"github.com/sirupsen/logrus"
	"sk_seckill_srv/conf"
	"streadway/streadway/amqp"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
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
	//确保每个消费者一次消费的数量，每次接受的消息数量
	err = ch.Qos(10, 0, false)
	if err != nil {
		logrus.WithField("err", err).Error("设置channel失败")
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

func (r *RabbitMQ) Init() *RabbitMQ {
	r.connMQ()
	r.openChan()
	return r
}

func (r *RabbitMQ) Close() {
	r.closeChan()
	r.closeConn()
}
