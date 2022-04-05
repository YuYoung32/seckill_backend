package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_product_srv "sk_product_srv/proto"
	. "sk_web/utils"
	"strconv"
)

type FrontProduct struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Num         int     `json:"num"`
	Unit        string  `json:"unit"`
	Pic         string  `json:"pic"`
	Description string  `json:"desc"`
	CreateTime  string  `json:"create_time"`
}

func GetProductsHandler(ctx *gin.Context) {
	var err error
	qcurrentPage := ctx.DefaultQuery("currentPage", "1")
	currentPage, err := strconv.Atoi(qcurrentPage)
	qpageSize := ctx.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(qpageSize)
	if err != nil {
		logrus.WithField("module", "get_user_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数错误",
		})
		return
	}
	start := (currentPage - 1) * pageSize

	service := grpc.NewService()
	client := yuyoung_srv_sk_product_srv.NewProductService("yuyoung.srv.sk_product_srv", service.Client())
	resp, err := client.GetProductList(context.TODO(), &yuyoung_srv_sk_product_srv.GetProductListRequest{
		Start:  int32(start),
		Amount: int32(pageSize),
	})
	if err != nil {
		logrus.WithField("module", "get_products_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    resp.GeneralResponse.Code,
			"message": resp.GeneralResponse.Msg,
		})
		return
	}

	var products []*FrontProduct
	//消息型数据转换成前端需要的数据，选择+更名json，因为消息类型数据是写死的，不可能随时根据业务调整
	for _, r := range resp.ProductInfo {
		products = append(products, &FrontProduct{
			Id:          r.ProductId,
			Name:        r.Name,
			Price:       r.Price,
			Num:         int(r.LeftNum),
			Unit:        r.Unit,
			Pic:         r.Image,
			Description: r.Description,
			CreateTime:  r.CreateTime,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":         resp.GeneralResponse.Code,
		"message":      resp.GeneralResponse.Msg,
		"products":     products,
		"page_size":    pageSize,
		"current_page": currentPage,
	})

}

func GetOneProductHandler(ctx *gin.Context) {
	productId := ctx.Query("id")

	service := grpc.NewService()
	client := yuyoung_srv_sk_product_srv.NewProductService("yuyoung.srv.sk_product_srv", service.Client())
	resp, err := client.GetProduct(context.TODO(), &yuyoung_srv_sk_product_srv.GeneralRequest{
		ProductId: productId,
	})
	if err != nil {
		logrus.WithField("module", "get_one_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusNotFound,
			"msg":  "商品不存在",
		})
		return
	}
	PicBase64, _ := Pic2Base64(resp.ProductInfo.Image)
	product := FrontProduct{
		Id:          resp.ProductInfo.ProductId,
		Name:        resp.ProductInfo.Name,
		Price:       resp.ProductInfo.Price,
		Num:         int(resp.ProductInfo.LeftNum),
		Unit:        resp.ProductInfo.Unit,
		Pic:         PicBase64,
		Description: resp.ProductInfo.Description,
		CreateTime:  resp.ProductInfo.CreateTime,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":       resp.GeneralResponse.Code,
		"msg":        resp.GeneralResponse.Msg,
		"product":    product,
		"img_base64": PicBase64,
	})

}
