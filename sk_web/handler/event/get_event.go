package event

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_product_srv "sk_product_srv/proto"
	yuyoung_srv_sk_event_srv "sk_product_srv/proto/event"
	"strconv"
)

type FrontEvent struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	ProductName string  `json:"pname"`
	EventPrice  float32 `json:"price"`
	EventNum    string  `json:"num"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	CreateTime  string  `json:"create_time"`
	ProductId   string  `json:"pid"`
}

func GetEventsHandler(ctx *gin.Context) {
	var err error
	qcurrentPage := ctx.DefaultQuery("currentPage", "1")
	currentPage, err := strconv.Atoi(qcurrentPage)
	qpageSize := ctx.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(qpageSize)
	if err != nil {
		logrus.WithField("module", "get_events_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
		})
		return
	}
	start := (currentPage - 1) * pageSize

	service := grpc.NewService()
	client := yuyoung_srv_sk_event_srv.NewEventService("yuyoung.srv.sk_event_srv", service.Client())
	resp, err := client.GetEventList(context.TODO(), &yuyoung_srv_sk_event_srv.GetEventListRequest{
		Start:  int32(start),
		Amount: int32(pageSize),
	})
	if err != nil {
		logrus.WithField("module", "get_events_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务错误",
		})
		return
	}

	var seckills []*FrontEvent
	//消息型数据转换成前端需要的数据，选择+更名json，因为消息类型数据是写死的，不可能随时根据业务调整
	for _, r := range resp.EventList {
		seckills = append(seckills, &FrontEvent{
			Id:          r.EventInfo.Id,
			Name:        r.EventInfo.Name,
			ProductName: r.ProductName,
			EventPrice:  r.EventInfo.EventPrice,
			EventNum:    strconv.Itoa(int(r.EventInfo.EventNum)),
			StartTime:   r.EventInfo.StartTime,
			EndTime:     r.EventInfo.EndTime,
			CreateTime:  r.EventCreateTime,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":         resp.GeneralResponse.Code,
		"msg":          resp.GeneralResponse.Msg,
		"seckills":     seckills,
		"total":        resp.Total,
		"page_size":    pageSize,
		"current_page": currentPage,
	})
}

func GetEventHandler(ctx *gin.Context) {
	eventId := ctx.Query("id")
	eventService := grpc.NewService()
	eventClient := yuyoung_srv_sk_event_srv.NewEventService("yuyoung.srv.sk_event_srv", eventService.Client())
	eventResp, err := eventClient.GetEvent(context.TODO(), &yuyoung_srv_sk_event_srv.GeneralRequest{
		EventId: eventId,
	})
	if err != nil {
		logrus.WithField("module", "get_event_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务错误",
		})
		return
	}

	productService := grpc.NewService()
	productClient := yuyoung_srv_sk_product_srv.NewProductService("yuyoung.srv.sk_product_srv", productService.Client())
	productResp, err := productClient.GetSelectedProductList(context.TODO(), &yuyoung_srv_sk_product_srv.GeneralRequest{
		ProductId: eventResp.EventDetailedInfo.EventInfo.ProductId,
	})
	if err != nil {
		logrus.WithField("module", "get_event_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": eventResp.GeneralResponse.Code,
			"msg":  eventResp.GeneralResponse.Msg,
		})
		return
	}

	type ProductsNo struct {
		ProductName string `json:"pname"`
		Id          string `json:"id"`
	}
	var productsNo []*ProductsNo
	for _, r := range productResp.ProductInfo {
		productsNo = append(productsNo, &ProductsNo{
			ProductName: r.Name,
			Id:          r.ProductId,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"seckill": &FrontEvent{
			Id:          eventResp.EventDetailedInfo.EventInfo.Id,
			Name:        eventResp.EventDetailedInfo.EventInfo.Name,
			ProductName: eventResp.EventDetailedInfo.ProductName,
			EventPrice:  eventResp.EventDetailedInfo.EventInfo.EventPrice,
			EventNum:    strconv.Itoa(int(eventResp.EventDetailedInfo.EventInfo.EventNum)),
			StartTime:   eventResp.EventDetailedInfo.EventInfo.StartTime,
			EndTime:     eventResp.EventDetailedInfo.EventInfo.EndTime,
			CreateTime:  eventResp.EventDetailedInfo.EventCreateTime,
			ProductId:   eventResp.EventDetailedInfo.EventInfo.ProductId,
		},
		"products_no": productsNo,
	})
}
