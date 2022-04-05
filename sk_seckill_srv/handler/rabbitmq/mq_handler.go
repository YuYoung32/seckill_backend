package rabbitmq

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"sk_seckill_srv/database"
	"sk_seckill_srv/handler/predis"
	. "sk_seckill_srv/utils"
	"strconv"
	"streadway/streadway/amqp"
	"time"
)

func GetSkMqHandler() {
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

func SkMqHandle(delivery amqp.Delivery) {
	redisConn := predis.GetRedisConn()
	defer redisConn.Close()
	db := database.GetDBConn()
	msg := StrToMap(string(delivery.Body))
	eventId := msg["eventId"].(string)
	email := msg["email"].(string)
	///*
	//	检验用户是否存在，并根据email查询Id
	// */
	var user database.User
	res := db.Where("email=?", email).First(&user)
	if res.Error != nil {
		logrus.WithField("module", "sk_mq_handler").Error(res.Error)
		_, err := redisConn.Do("set", email, MapToStr(map[string]interface{}{"code": http.StatusInternalServerError, "msg": "参数错误"}))
		if err != nil {
			logrus.WithField("module", "sk_mq_handler").Error(err)
		}
		delivery.Ack(false)
		return
	}
	userId := user.ID

	/*
		订单验证
		1.时间限制，在规定的时间内才能进行购买
		2.活动验证，活动是否存在
		3.数量限制，商品数量等于0，则不能进行购买
		4.参与限制，一个用户只能参与一个活动
	*/
	nowTime := time.Now()
	type ProductSelling struct {
		ProductId   int
		ProductName string
		LeftNum     int
	}
	var productSelling ProductSelling
	//查询判断是否是符合时间的活动，并查询商品信息
	res = db.Table("events").Select("events.product_id as product_id, products.name as product_name, products.left_num as left_num").
		Joins("left join products on events.product_id = products.id").
		Where("events.start_time <= ?", nowTime).Where("events.end_time >= ?", nowTime).Where("events.id = ?", eventId).
		Scan(&productSelling)
	if res.Error != nil || res.RowsAffected != 1 || productSelling.LeftNum <= 0 {
		if res.Error != nil {
			logrus.WithField("module", "sk_mq_handler").Error(res.Error)
			delivery.Ack(false)
			return
		}
		logrus.WithField("userId", userId).Debug("库存不足或已参与过活动或者不在活动时间内")
		_, err := redisConn.Do("set", email, MapToStr(map[string]interface{}{"code": http.StatusNotFound, "msg": "库存不足或已参与过活动或者不在活动时间内"}))
		if err != nil {
			logrus.WithField("module", "sk_mq_handler").Error(err)
		}
		delivery.Ack(false)
		return
	}
	//查询是否已经参与过活动
	res = db.Model(&database.Order{}).Where("user_id=? and event_id=?", userId, eventId).Scan(&database.Order{})
	if res.RowsAffected != 0 {
		logrus.WithField("module", "sk_mq_handler").WithField("userId", userId).Debug("已经参与过活动")
		_, err := redisConn.Do("set", email, MapToStr(map[string]interface{}{"code": http.StatusNotFound, "msg": "已经参与过该活动"}))
		if err != nil {
			logrus.WithField("module", "sk_mq_handler").Error(err)
		}
		delivery.Ack(false)
		return
	}

	/*
		订单处理
		1.库存减一
		2.生成订单
	*/
	//商品表更新
	res = db.Model(database.Product{}).Where("id = ?", productSelling.ProductId).Update("left_num", productSelling.LeftNum-1)
	var event database.Event
	//活动表更新
	res = db.Model(database.Event{}).Where("id = ?", eventId).Find(&event)
	res = db.Model(database.Event{}).Where("id = ?", eventId).Update("event_num", event.EventNum-1)
	if res.Error != nil {
		logrus.WithField("module", "sk_mq_handler").Error(res.Error)
		_, err := redisConn.Do("set", email, MapToStr(map[string]interface{}{"code": http.StatusInternalServerError, "msg": "系统错误"}))
		if err != nil {
			logrus.WithField("module", "sk_mq_handler").Error(err)
		}
		delivery.Ack(false)
		return
	}

	order := new(database.Order)
	atoi, err := strconv.Atoi(eventId)
	if err != nil {
		logrus.WithField("module", "sk_mq_handler").Error(err)
		_, err := redisConn.Do("set", email, MapToStr(map[string]interface{}{"code": http.StatusNotFound, "msg": "参数错误"}))
		if err != nil {
			logrus.WithField("module", "sk_mq_handler").Error(err)
		}
		delivery.Ack(false)
		return
	}
	order.OrderSerial = strconv.Itoa(int(time.Now().UnixNano())) + "E" + strconv.Itoa(atoi) + "P" + strconv.Itoa(productSelling.ProductId)
	order.EventId = uint(atoi)
	order.UserId = userId
	order.PayStatus = "0"
	res = db.Create(&order)
	if res.Error != nil {
		logrus.WithField("module", "sk_mq_handler").Error(res.Error)
		_, err := redisConn.Do("set", email, MapToStr(map[string]interface{}{"code": http.StatusInternalServerError, "msg": "订单生成失败"}))
		if err != nil {
			logrus.WithField("module", "sk_mq_handler").Error(err)
		}
		delivery.Ack(false)
		return
	}

	logrus.WithField("module", "sk_mq_handler").WithFields(logrus.Fields{
		"event_id":   eventId,
		"user_id":    userId,
		"product_id": productSelling.ProductId,
	}).Info("订单生成成功")
	_, err = redisConn.Do("set", email, MapToStr(map[string]interface{}{"code": http.StatusOK, "msg": "订单生成成功"}))
	if err != nil {
		logrus.WithField("module", "sk_mq_handler").Error(err)
	}
	delivery.Ack(false)
	return
}
