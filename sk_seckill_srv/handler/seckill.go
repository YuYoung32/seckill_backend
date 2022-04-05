package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"sk_seckill_srv/database"
	seckill "sk_seckill_srv/proto"
	"strconv"
	"time"
)

type SeckillImpl struct {
}

func (s SeckillImpl) FrontSeckill(ctx context.Context, in *seckill.SeckillRequest, out *seckill.GeneralResponse) error {
	db := database.GetDBConn()
	eventId := in.Id
	var user database.User
	res := db.Where("email=?", in.Email).First(&user)
	userId := user.ID
	if res.Error != nil {
		logrus.WithField("module", "front_seckill").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "参数错误"
		return nil
	}

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
	res = db.Table("events").Select("events.product_id as product_id, products.name as product_name, products.left_num as left_num").
		Joins("left join products on events.product_id = products.id").
		Where("events.start_time <= ?", nowTime).Where("events.end_time >= ?", nowTime).Where("events.id = ?", eventId).
		Scan(&productSelling)
	if res.Error != nil || res.RowsAffected != 1 || productSelling.LeftNum <= 0 {
		logrus.WithField("module", "front_seckill").Info(res.Error)
		out.Code = strconv.Itoa(http.StatusNotFound)
		out.Msg = "已经售罄或不在指定时间段或活动不存在"
		return nil
	}
	res = db.Model(&database.Order{}).Where("user_id=? and event_id=?", userId, eventId).Scan(&database.Order{})
	if res.RowsAffected != 0 {
		logrus.WithField("module", "front_seckill").WithField("userId", userId).Info(res.Error)
		out.Code = strconv.Itoa(http.StatusNotFound)
		out.Msg = "已经参与过该活动"
		return nil
	}

	/*
		订单处理
		1.库存减一
		2.生成订单
	*/
	res = db.Model(database.Product{}).Where("id = ?", productSelling.ProductId).Update("left_num", productSelling.LeftNum-1)
	var event database.Event
	res = db.Model(database.Event{}).Where("id = ?", eventId).Find(&event)
	res = db.Model(database.Event{}).Where("id = ?", eventId).Update("event_num", event.EventNum-1)
	if res.Error != nil {
		logrus.WithField("module", "front_seckill").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "系统错误"
		return nil
	}

	order := new(database.Order)
	atoi, err := strconv.Atoi(eventId)
	if err != nil {
		logrus.WithField("module", "front_seckill").Error(err)
		out.Code = strconv.Itoa(http.StatusNotFound)
		out.Msg = "参数错误"
		return nil
	}
	order.OrderSerial = strconv.Itoa(int(time.Now().UnixNano())) + "E" + strconv.Itoa(atoi) + "P" + strconv.Itoa(productSelling.ProductId)
	order.EventId = uint(atoi)
	order.UserId = userId
	order.PayStatus = "0"
	res = db.Create(&order)
	if res.Error != nil {
		logrus.WithField("module", "front_seckill").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "订单生成失败"
		return nil
	}

	logrus.WithField("module", "front_seckill").WithFields(logrus.Fields{
		"event_id":   eventId,
		"user_id":    userId,
		"product_id": productSelling.ProductId,
	}).Info("订单生成成功")
	out.Code = strconv.Itoa(http.StatusOK)
	out.Msg = "订单生成成功"
	return nil

}
