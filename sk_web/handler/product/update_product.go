package product

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	yuyoung_srv_sk_product_srv "sk_product_srv/proto"
	"strconv"
	"time"
)

func UpdateProductHandler(ctx *gin.Context) {
	var err error
	id := ctx.PostForm("id")
	name := ctx.PostForm("name")
	pprice := ctx.PostForm("price")
	price, err := strconv.ParseFloat(pprice, 32)
	pnum := ctx.PostForm("num")
	num, err := strconv.Atoi(pnum)
	unit := ctx.PostForm("unit")
	description := ctx.PostForm("description")
	picture, err := ctx.FormFile("pic")
	if err != nil {
		logrus.WithField("module", "update_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "信息错误",
		})
		return
	}

	filePath := "handler/product/product_pic/" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + picture.Filename

	service := grpc.NewService()
	client := yuyoung_srv_sk_product_srv.NewProductService("yuyoung.srv.sk_product_srv", service.Client())
	//查找原先的图片
	pic_resp, err := client.GetProduct(context.TODO(), &yuyoung_srv_sk_product_srv.GeneralRequest{
		ProductId: id,
	})
	if err != nil {
		logrus.WithField("module", "update_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": pic_resp.GeneralResponse.Code,
			"msg":  pic_resp.GeneralResponse.Msg,
		})
		return
	}
	toDelete := pic_resp.ProductInfo.Image

	resp, err := client.EditProduct(context.TODO(), &yuyoung_srv_sk_product_srv.EditProductRequest{
		ProductInfo: &yuyoung_srv_sk_product_srv.ProductInfo{
			ProductId:   id,
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
		logrus.WithField("module", "update_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	err = ctx.SaveUploadedFile(picture, filePath)
	if err != nil {
		logrus.WithField("module", "update_product_handler").Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "保存图片失败",
		})
		return
	}
	_ = os.Remove(toDelete)

	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}
