package event

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"net/http"
	yuyoung_srv_sk_product_srv "sk_product_srv/proto"
)

func GetProductsHandler(ctx *gin.Context) {
	service := grpc.NewService()
	client := yuyoung_srv_sk_product_srv.NewProductService("yuyoung.srv.sk_product_srv", service.Client())
	resp, err := client.GetProductList(context.TODO(), &yuyoung_srv_sk_product_srv.GetProductListRequest{
		Start:  0,
		Amount: -1,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.GeneralResponse.Code,
			"msg":  resp.GeneralResponse.Msg,
		})
		return
	}
	type FrontProduct struct {
		Id   string `json:"id"`
		Name string `json:"pname"`
	}
	var products []FrontProduct
	for _, r := range resp.ProductInfo {
		products = append(products, FrontProduct{
			Id:   r.ProductId,
			Name: r.Name,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":     resp.GeneralResponse.Code,
		"msg":      resp.GeneralResponse.Msg,
		"products": products,
	})
}
