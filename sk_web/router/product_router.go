package router

import (
	"github.com/gin-gonic/gin"
	"sk_web/handler/product"
	"sk_web/middleware"
)

func ProductRouter(router *gin.RouterGroup) {
	router.GET("/get_products", middleware.AdminAuth, product.GetProductsHandler)
	router.POST("/add_product", middleware.AdminAuth, product.AddProductHandler)
	router.GET("/get_one_product", middleware.AdminAuth, product.GetOneProductHandler)
	router.POST("/update_product", middleware.AdminAuth, product.UpdateProductHandler)
	router.POST("/delete_product", middleware.AdminAuth, product.DeleteProductHandler)

}
