package event

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_event_srv "sk_product_srv/proto/event"
	. "sk_web/utils"
	"strconv"
)

type FrontEventWithProduct struct {
	EventId    string  `json:"id"` //活动ID
	LeftNum    int     `json:"num"`
	EventPrice float32 `json:"price"` //活动价格
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`

	ProductId          string  `json:"pid"`
	ProductName        string  `json:"name"`    //商品名称
	Picture            string  `json:"pic"`     //base64
	Price              float32 `json:"p_price"` //原价
	ProductDescription string  `json:"pdesc"`
	Unit               string  `json:"unit"`
}

func GetFrontEventsHandler(ctx *gin.Context) {
	var err error
	qcurrentPage := ctx.DefaultQuery("currentPage", "1")
	currentPage, err := strconv.Atoi(qcurrentPage)
	qpageSize := ctx.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(qpageSize)
	if err != nil {
		logrus.WithField("module", "get_front_events_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
		})
		return
	}
	start := (currentPage - 1) * pageSize

	service := grpc.NewService()
	client := yuyoung_srv_sk_event_srv.NewEventService("yuyoung.srv.sk_event_srv", service.Client())
	resp, err := client.GetFrontEventList(context.TODO(), &yuyoung_srv_sk_event_srv.GetFrontEventListRequest{
		Start:  int32(start),
		Amount: int32(pageSize),
	})
	if err != nil {
		logrus.WithField("module", "get_front_events_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务错误",
		})
		return
	}

	var seckill_list []*FrontEventWithProduct
	//消息型数据转换成前端需要的数据，选择+更名json，因为消息类型数据是写死的，不可能随时根据业务调整
	for _, r := range resp.FrontEventList {
		picBase64, _ := Pic2Base64(r.ProductInfo.Image)
		seckill_list = append(seckill_list, &FrontEventWithProduct{
			EventId:            r.EventDetailedInfo.EventInfo.Id,
			LeftNum:            int(r.EventDetailedInfo.EventInfo.EventNum),
			EventPrice:         r.EventDetailedInfo.EventInfo.EventPrice,
			StartTime:          r.EventDetailedInfo.EventInfo.StartTime,
			EndTime:            r.EventDetailedInfo.EventInfo.EndTime,
			ProductId:          r.ProductInfo.ProductId,
			ProductName:        r.ProductInfo.Name,
			Picture:            picBase64,
			Price:              r.ProductInfo.Price,
			ProductDescription: r.ProductInfo.Description,
			Unit:               r.ProductInfo.Unit,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":         resp.GeneralResponse.Code,
		"msg":          resp.GeneralResponse.Msg,
		"seckill_list": seckill_list,
		"total_page":   int(resp.Total) / pageSize,
		"page_size":    pageSize,
		"current_page": currentPage,
	})
}

func GetFrontEventHandler(ctx *gin.Context) {
	eventId := ctx.Query("id")
	eventService := grpc.NewService()
	eventClient := yuyoung_srv_sk_event_srv.NewEventService("yuyoung.srv.sk_event_srv", eventService.Client())
	eventResp, err := eventClient.GetFrontEvent(context.TODO(), &yuyoung_srv_sk_event_srv.GeneralRequest{
		EventId: eventId,
	})
	if err != nil {
		logrus.WithField("module", "get_front_event_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务错误",
		})
		return
	}

	picBase64, _ := Pic2Base64(eventResp.FrontEventInfo.ProductInfo.Image)

	ctx.JSON(http.StatusOK, gin.H{
		"code": eventResp.GeneralResponse.Code,
		"msg":  eventResp.GeneralResponse.Msg,
		"seckill": FrontEventWithProduct{
			EventId:            eventResp.FrontEventInfo.EventDetailedInfo.EventInfo.Id,
			LeftNum:            int(eventResp.FrontEventInfo.EventDetailedInfo.EventInfo.EventNum),
			EventPrice:         eventResp.FrontEventInfo.EventDetailedInfo.EventInfo.EventPrice,
			StartTime:          eventResp.FrontEventInfo.EventDetailedInfo.EventInfo.StartTime,
			EndTime:            eventResp.FrontEventInfo.EventDetailedInfo.EventInfo.EndTime,
			ProductId:          eventResp.FrontEventInfo.ProductInfo.ProductId,
			ProductName:        eventResp.FrontEventInfo.ProductInfo.Name,
			Picture:            picBase64,
			Price:              eventResp.FrontEventInfo.ProductInfo.Price,
			ProductDescription: eventResp.FrontEventInfo.ProductInfo.Description,
			Unit:               eventResp.FrontEventInfo.ProductInfo.Unit,
		},
	})
}
