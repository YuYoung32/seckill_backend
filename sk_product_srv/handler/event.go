package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"sk_product_srv/database"
	yuyoung_srv_sk_event_srv "sk_product_srv/proto/event"
	"strconv"
	"time"
)

type EventImpl struct {
}

type EventDetailedInfo struct {
	EventId         string
	EventCreateTime time.Time
	Name            string
	EventPrice      float32
	EventNum        int
	StartTime       time.Time
	EndTime         time.Time
	ProductId       uint
	ProductName     string
}

type EventDetailedWithProduct struct {
	EventId    string //活动ID
	LeftNum    int
	EventPrice float32 //活动价格
	StartTime  string
	EndTime    string

	ProductId          string
	ProductName        string  //商品名称
	Picture            string  //base64
	Price              float32 //原价
	ProductDescription string
	Unit               string
}

func (e EventImpl) AddEvent(ctx context.Context, in *yuyoung_srv_sk_event_srv.AddEventRequest, out *yuyoung_srv_sk_event_srv.GeneralResponse) error {
	event := new(database.Event)
	event.Name = in.EventInfo.Name
	event.EventPrice = in.EventInfo.EventPrice
	event.EventNum = int(in.EventInfo.EventNum)

	startTime, err := time.Parse("2006-01-02 15:04:05", in.EventInfo.StartTime)
	if err != nil {
		logrus.WithField("module", "add_event").Error(err)
		out.Msg = err.Error()
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		return err
	}
	event.StartTime = startTime

	endTime, err := time.Parse("2006-01-02 15:04:05", in.EventInfo.EndTime)
	if err != nil {
		logrus.WithField("module", "add_event").Error(err)
		out.Msg = err.Error()
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		return err
	}
	event.EndTime = endTime

	productId, err := strconv.Atoi(in.EventInfo.ProductId)
	if err != nil {
		logrus.WithField("module", "add_event").Error(err)
		out.Msg = err.Error()
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		return err
	}
	event.ProductId = uint(productId)

	db := database.GetDBConn()
	res := db.Create(&event)
	if res.Error != nil {
		logrus.WithField("module", "add_event").Error(res.Error)
		out.Msg = res.Error.Error()
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		return res.Error
	}
	logrus.WithField("module", "add_event").WithField("活动名称", event.Name).Info("更新活动成功")
	out.Msg = "添加活动成功"
	out.Code = strconv.Itoa(http.StatusOK)
	return nil
}

func (e EventImpl) GetEventList(ctx context.Context, in *yuyoung_srv_sk_event_srv.GetEventListRequest, out *yuyoung_srv_sk_event_srv.GetEventListResponse) error {
	db := database.GetDBConn()
	var events []EventDetailedInfo
	var total int
	res := db.Table("events").Select("events.id as event_id, events.created_at as event_create_time, events.name, events.event_price, events.event_num, events.start_time ,events.end_time, events.product_id, products.name as product_name").
		Joins("left join products on events.product_id = products.id").Offset(in.Start).Limit(in.Amount).Scan(&events)
	res = db.Model(&database.Event{}).Count(&total)
	if res.Error != nil {
		logrus.WithField("module", "get_event_list").Error(res.Error)
		out.GeneralResponse.Code = strconv.Itoa(http.StatusInternalServerError)
		out.GeneralResponse.Msg = res.Error.Error()
		return res.Error
	}

	var frontEvents []*yuyoung_srv_sk_event_srv.EventDetailedInfo
	for _, e := range events {
		frontEvents = append(frontEvents, &yuyoung_srv_sk_event_srv.EventDetailedInfo{
			EventInfo: &yuyoung_srv_sk_event_srv.EventInfo{
				Id:         e.EventId,
				Name:       e.Name,
				EventPrice: e.EventPrice,
				EventNum:   int32(e.EventNum),
				StartTime:  e.StartTime.String(),
				EndTime:    e.EndTime.String(),
				ProductId:  strconv.Itoa(int(e.ProductId)),
			},
			EventCreateTime: e.EventCreateTime.String(),
			ProductName:     e.ProductName,
		})
	}

	*out = yuyoung_srv_sk_event_srv.GetEventListResponse{
		GeneralResponse: &yuyoung_srv_sk_event_srv.GeneralResponse{
			Code: strconv.Itoa(http.StatusOK),
			Msg:  "获取活动列表成功",
		},
		EventList: frontEvents,
		Total:     int32(total),
	}
	return nil
}

func (e EventImpl) GetEvent(ctx context.Context, in *yuyoung_srv_sk_event_srv.GeneralRequest, out *yuyoung_srv_sk_event_srv.GetEventResponse) error {
	db := database.GetDBConn()

	var event EventDetailedInfo
	res := db.Table("events").Where("events.id=?", in.EventId).Select("events.id as event_id, events.created_at as event_create_time, events.name, events.event_price, events.event_num, events.start_time ,events.end_time, events.product_id, products.name as product_name").
		Joins("left join products on events.product_id = products.id").Scan(&event)
	if res.Error != nil {
		logrus.WithField("module", "get_event").Error(res.Error)
		out.GeneralResponse.Code = strconv.Itoa(http.StatusInternalServerError)
		out.GeneralResponse.Msg = res.Error.Error()
		return res.Error
	}

	*out = yuyoung_srv_sk_event_srv.GetEventResponse{
		GeneralResponse: &yuyoung_srv_sk_event_srv.GeneralResponse{
			Code: strconv.Itoa(http.StatusOK),
			Msg:  "获取活动详情成功",
		},
		EventDetailedInfo: &yuyoung_srv_sk_event_srv.EventDetailedInfo{
			EventInfo: &yuyoung_srv_sk_event_srv.EventInfo{
				Id:         event.EventId,
				Name:       event.Name,
				EventPrice: event.EventPrice,
				EventNum:   int32(event.EventNum),
				StartTime:  event.StartTime.String(),
				EndTime:    event.EndTime.String(),
				ProductId:  strconv.Itoa(int(event.ProductId)),
			},
			EventCreateTime: event.EventCreateTime.String(),
			ProductName:     event.ProductName,
		},
	}
	return nil
}

func (e EventImpl) EditEvent(ctx context.Context, in *yuyoung_srv_sk_event_srv.EditEventRequest, out *yuyoung_srv_sk_event_srv.GeneralResponse) error {
	db := database.GetDBConn()
	var event database.Event

	event.Name = in.EventInfo.Name
	event.EventPrice = in.EventInfo.EventPrice
	event.EventNum = int(in.EventInfo.EventNum)
	startTime, err := time.Parse("2006-01-02 15:04:05", in.EventInfo.StartTime)
	endTime, err := time.Parse("2006-01-02 15:04:05", in.EventInfo.EndTime)
	if err != nil {
		logrus.WithField("module", "edit_event").Error(err)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "传入日期错误"
		return err
	}
	event.StartTime = startTime
	event.EndTime = endTime
	productId, err := strconv.Atoi(in.EventInfo.ProductId)
	if err != nil {
		logrus.WithField("module", "edit_event").Error(err)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "传入产品id错误"
		return err
	}
	event.ProductId = uint(productId)

	res := db.Model(&database.Event{}).Where("id=?", in.EventInfo.Id).Update(&event)
	if res.Error != nil {
		logrus.WithField("module", "edit_event").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = res.Error.Error()
		return res.Error
	}
	logrus.WithField("module", "edit_event").WithField("活动id", in.EventInfo.Id).WithField("活动名称", event.Name).Info("修改活动成功")
	out.Code = strconv.Itoa(http.StatusOK)
	out.Msg = "修改活动成功"
	return nil
}

func (e EventImpl) DeleteEvent(ctx context.Context, in *yuyoung_srv_sk_event_srv.GeneralRequest, out *yuyoung_srv_sk_event_srv.GeneralResponse) error {
	db := database.GetDBConn()
	res := db.Delete(&database.Event{}, "id=?", in.EventId)
	if res.Error != nil {
		logrus.WithField("module", "delete_event").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = res.Error.Error()
		return res.Error
	}
	out.Msg = "删除活动成功"
	out.Code = strconv.Itoa(http.StatusOK)
	logrus.WithField("module", "delete_event").WithField("活动id", in.EventId).Info("删除活动成功")
	return nil
}

func (e EventImpl) GetFrontEventList(ctx context.Context, in *yuyoung_srv_sk_event_srv.GetFrontEventListRequest, out *yuyoung_srv_sk_event_srv.GetFrontEventListResponse) error {
	db := database.GetDBConn()
	var events []EventDetailedWithProduct
	var total int

	res := db.Table("events").Joins("left join products on events.product_id = products.id").
		Select("events.id            as event_id," +
			"       events.event_num     as left_num," +
			"       events.event_price   as event_price," +
			"       events.start_time    as start_time," +
			"       events.end_time      as end_time," +
			"       products.id          as product_id," +
			"       products.name        as product_name," +
			"       products.picture     as picture," +
			"       products.price       as price," +
			"       products.description as product_description," +
			"       products.unit        as unit").Offset(in.Start).Limit(in.Amount).
		Scan(&events)
	res = db.Model(&database.Event{}).Count(&total)
	if res.Error != nil {
		logrus.WithField("module", "get_front_events_list").Error(res.Error)
		out.GeneralResponse.Code = strconv.Itoa(http.StatusInternalServerError)
		out.GeneralResponse.Msg = res.Error.Error()
		return res.Error
	}

	var frontEvents []*yuyoung_srv_sk_event_srv.FrontEventInfo
	for _, e := range events {
		frontEvents = append(frontEvents, &yuyoung_srv_sk_event_srv.FrontEventInfo{
			EventDetailedInfo: &yuyoung_srv_sk_event_srv.EventDetailedInfo{
				EventInfo: &yuyoung_srv_sk_event_srv.EventInfo{
					Id:         e.EventId,
					EventPrice: e.EventPrice,
					EventNum:   int32(e.LeftNum),
					StartTime:  e.StartTime,
					EndTime:    e.EndTime,
				},
			},
			ProductInfo: &yuyoung_srv_sk_event_srv.ProductInfo{
				ProductId:   e.ProductId,
				Name:        e.ProductName,
				Price:       e.Price,
				Unit:        e.Unit,
				Image:       e.Picture,
				Description: e.ProductDescription,
			},
		})
	}

	*out = yuyoung_srv_sk_event_srv.GetFrontEventListResponse{
		GeneralResponse: &yuyoung_srv_sk_event_srv.GeneralResponse{
			Code: strconv.Itoa(http.StatusOK),
			Msg:  "查询成功",
		},
		Total:          int32(total),
		FrontEventList: frontEvents,
	}
	return nil
}

func (e EventImpl) GetFrontEvent(ctx context.Context, in *yuyoung_srv_sk_event_srv.GeneralRequest, out *yuyoung_srv_sk_event_srv.GetFrontEventResponse) error {
	db := database.GetDBConn()
	var event EventDetailedWithProduct
	var total int

	res := db.Table("events").Joins("left join products on events.product_id = products.id").
		Select("events.id            as event_id,"+
			"       events.event_num     as left_num,"+
			"       events.event_price   as event_price,"+
			"       events.start_time    as start_time,"+
			"       events.end_time      as end_time,"+
			"       products.id          as product_id,"+
			"       products.name        as product_name,"+
			"       products.picture     as picture,"+
			"       products.price       as price,"+
			"       products.description as product_description,"+
			"       products.unit        as unit").Where("events.id=?", in.EventId).Scan(&event)
	res = db.Model(&database.Event{}).Count(&total)
	if res.Error != nil {
		logrus.WithField("module", "get_front_event_list").Error(res.Error)
		out.GeneralResponse.Code = strconv.Itoa(http.StatusInternalServerError)
		out.GeneralResponse.Msg = res.Error.Error()
		return res.Error
	}

	*out = yuyoung_srv_sk_event_srv.GetFrontEventResponse{
		GeneralResponse: &yuyoung_srv_sk_event_srv.GeneralResponse{
			Code: strconv.Itoa(http.StatusOK),
			Msg:  "查询成功",
		},
		FrontEventInfo: &yuyoung_srv_sk_event_srv.FrontEventInfo{
			EventDetailedInfo: &yuyoung_srv_sk_event_srv.EventDetailedInfo{
				EventInfo: &yuyoung_srv_sk_event_srv.EventInfo{
					Id:         event.EventId,
					EventPrice: event.EventPrice,
					EventNum:   int32(event.LeftNum),
					StartTime:  event.StartTime,
					EndTime:    event.EndTime,
				},
			},
			ProductInfo: &yuyoung_srv_sk_event_srv.ProductInfo{
				ProductId:   event.ProductId,
				Name:        event.ProductName,
				Price:       event.Price,
				Unit:        event.Unit,
				Image:       event.Picture,
				Description: event.ProductDescription,
			},
		},
	}
	return nil

}
