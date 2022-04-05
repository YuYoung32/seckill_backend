package handler

import (
	. "context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	"sk_product_srv/database"
	. "sk_product_srv/proto"
	"strconv"
)

type ProductImpl struct {
}

func (p ProductImpl) AddProduct(ctx Context, in *AddProductRequest, out *GeneralResponse) error {
	db := database.GetDBConn()
	product := new(database.Product)
	product.Name = in.ProductInfo.Name
	product.Price = in.ProductInfo.Price
	product.LeftNum = int(in.ProductInfo.LeftNum)
	product.Unit = in.ProductInfo.Unit
	product.Picture = in.ProductInfo.Image
	product.Description = in.ProductInfo.Description
	res := db.Create(&product)
	if res.Error != nil {
		logrus.WithField("module", "add_product").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "添加商品失败"
		return res.Error
	}

	logrus.WithField("商品名称", product.Name).Info("添加成功")
	out.Code = strconv.Itoa(http.StatusOK)
	out.Msg = "添加商品成功"
	return nil
}

func (p ProductImpl) DeleteProduct(ctx Context, in *GeneralRequest, out *GeneralResponse) error {
	db := database.GetDBConn()
	res := db.Model(database.Product{}).Where("id = ?", in.ProductId).Delete(&database.Product{})
	if res.Error != nil {
		logrus.WithField("module", "delete_product").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "删除商品失败"
		return res.Error
	}

	logrus.WithField("商品id", in.ProductId).Info("删除成功")
	out.Code = strconv.Itoa(http.StatusOK)
	out.Msg = "删除商品成功"
	return nil
}

func (p ProductImpl) GetProductList(ctx Context, in *GetProductListRequest, out *GetProductListResponse) error {
	db := database.GetDBConn()
	start := in.Start
	amount := in.Amount

	var res []database.Product
	var e *gorm.DB
	var total int
	if amount == -1 {
		e = db.Model(database.Product{}).Find(&res)
	} else {
		e = db.Limit(amount).Offset(start).Find(&res)
		e = db.Model(&database.Product{}).Count(&total)
	}
	if e.Error != nil {
		logrus.WithField("module", "get_product_list").Error(e.Error)
		out.GeneralResponse = &GeneralResponse{
			Code: strconv.Itoa(http.StatusInternalServerError),
			Msg:  "数据库查询错误",
		}
		return e.Error
	}

	//数据库查出数据转换成proto消息类型数据，选择+更名，因为消息类型数据是写死的，不可能和数据库字段一致
	var productInfo []*ProductInfo
	for _, r := range res {
		productInfo = append(productInfo, &ProductInfo{
			ProductId:   strconv.Itoa(int(r.ID)),
			Name:        r.Name,
			Price:       r.Price,
			LeftNum:     int32(r.LeftNum),
			Unit:        r.Unit,
			Image:       r.Picture,
			Description: r.Description,
			CreateTime:  r.CreatedAt.String(),
		})
	}

	out.GeneralResponse = &GeneralResponse{
		Code: strconv.Itoa(http.StatusOK),
		Msg:  "查询成功",
	}
	out.ProductInfo = productInfo
	out.Total = int32(total)

	return nil
}

func (p ProductImpl) EditProduct(ctx Context, in *EditProductRequest, out *GeneralResponse) error {
	db := database.GetDBConn()
	product := new(database.Product)
	product.Name = in.ProductInfo.Name
	product.Price = in.ProductInfo.Price
	product.LeftNum = int(in.ProductInfo.LeftNum)
	product.Unit = in.ProductInfo.Unit
	product.Picture = in.ProductInfo.Image
	product.Description = in.ProductInfo.Description
	res := db.Model(&database.Product{}).Where("id=?", in.ProductInfo.ProductId).Update(&product)
	if res.Error != nil {
		logrus.WithField("module", "edit_product").Error(res.Error)
		out.Code = strconv.Itoa(http.StatusInternalServerError)
		out.Msg = "更新商品失败"
		return res.Error
	}
	logrus.WithField("商品名称", product.Name).WithField("商品id", in.ProductInfo.ProductId).Info("更新成功")
	out.Code = strconv.Itoa(http.StatusOK)
	out.Msg = "更新商品成功"
	return nil
}

func (p ProductImpl) GetProduct(ctx Context, in *GeneralRequest, out *GetProductResponse) error {
	db := database.GetDBConn()
	var product database.Product
	res := db.Where("id=?", in.ProductId).First(&product)
	if res.Error != nil {
		logrus.WithField("module", "get_product").Error(res.Error)
		out.GeneralResponse = &GeneralResponse{
			Code: strconv.Itoa(http.StatusInternalServerError),
			Msg:  "查询失败",
		}
		return res.Error
	}
	if res.RowsAffected != 1 {
		logrus.WithField("module", "get_product").Error("多个查询结果")
		out.GeneralResponse = &GeneralResponse{
			Code: strconv.Itoa(http.StatusInternalServerError),
			Msg:  "多个查询结果",
		}
		return errors.New("多个查询结果")
	}
	out.GeneralResponse = &GeneralResponse{
		Code: strconv.Itoa(http.StatusOK),
		Msg:  "查询成功",
	}
	out.ProductInfo = &ProductInfo{
		ProductId:   strconv.Itoa(int(product.ID)),
		Name:        product.Name,
		Price:       product.Price,
		LeftNum:     int32(product.LeftNum),
		Unit:        product.Unit,
		Image:       product.Picture,
		Description: product.Description,
		CreateTime:  product.CreatedAt.String(),
	}
	return nil
}

func (p ProductImpl) GetSelectedProductList(ctx Context, in *GeneralRequest, out *GetProductListResponse) error {
	db := database.GetDBConn()
	var products []database.Product
	res := db.Model(database.Product{}).Where("id != ?", in.ProductId).Find(&products)
	if res.Error != nil {
		logrus.WithField("module", "get_selected_product_list").Error(res.Error)
		out.GeneralResponse = &GeneralResponse{
			Code: strconv.Itoa(http.StatusInternalServerError),
			Msg:  "查询失败",
		}
		return res.Error
	}
	var productInfo []*ProductInfo
	for _, r := range products {
		productInfo = append(productInfo, &ProductInfo{
			ProductId: strconv.Itoa(int(r.ID)),
			Name:      r.Name,
		})
	}
	*out = GetProductListResponse{
		GeneralResponse: &GeneralResponse{
			Code: strconv.Itoa(http.StatusOK),
			Msg:  "查询成功",
		},
		ProductInfo: productInfo,
	}
	return nil
}
