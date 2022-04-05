package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	yuyoung_srv_sk_product_srv "sk_product_srv/proto"
	"strconv"
	"time"
)

func AddProductHandler(ctx *gin.Context) {
	var err error
	name := ctx.PostForm("name")
	pprice := ctx.PostForm("price")
	price, err := strconv.ParseFloat(pprice, 32)
	pnum := ctx.PostForm("num")
	num, err := strconv.Atoi(pnum)
	unit := ctx.PostForm("unit")
	description := ctx.PostForm("description")
	picture, err := ctx.FormFile("pic")
	if err != nil {
		logrus.WithField("module", "add_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "信息错误",
		})
		return
	}

	filePath := "handler/product/product_pic/" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + picture.Filename

	service := grpc.NewService()
	client := yuyoung_srv_sk_product_srv.NewProductService("yuyoung.srv.sk_product_srv", service.Client())
	resp, err := client.AddProduct(context.TODO(), &yuyoung_srv_sk_product_srv.AddProductRequest{
		ProductInfo: &yuyoung_srv_sk_product_srv.ProductInfo{
			ProductId:   "", //不需要也不能提供
			Name:        name,
			Price:       float32(price),
			LeftNum:     int32(num),
			Unit:        unit,
			Image:       filePath,
			Description: description,
			CreateTime:  "", //不需要提供
		},
	})
	if err != nil {
		logrus.WithField("module", "add_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	err = ctx.SaveUploadedFile(picture, filePath)
	if err != nil {
		logrus.WithField("module", "add_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "保存图片失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}
