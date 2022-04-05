package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	yuyoung_srv_sk_product_srv "sk_product_srv/proto"
)

func DeleteProductHandler(ctx *gin.Context) {
	productId := ctx.PostForm("id")
	service := grpc.NewService()
	client := yuyoung_srv_sk_product_srv.NewProductService("yuyoung.srv.sk_product_srv", service.Client())

	//查找原先的图片
	pic_resp, err := client.GetProduct(context.TODO(), &yuyoung_srv_sk_product_srv.GeneralRequest{
		ProductId: productId,
	})
	if err != nil {
		logrus.WithField("module", "delete_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "商品不存在",
		})
		return
	}
	toDelete := pic_resp.ProductInfo.Image

	resp, err := client.DeleteProduct(context.TODO(), &yuyoung_srv_sk_product_srv.GeneralRequest{
		ProductId: productId,
	})
	if err != nil {
		logrus.WithField("module", "delete_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	_ = os.Remove(toDelete)

	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}
