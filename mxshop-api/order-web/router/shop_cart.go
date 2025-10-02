package router

import (
	shopcart "mxshop-api/order-web/api/shop_cart"
	"mxshop-api/order-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		ShopCartRouter.GET("", shopcart.List)
		ShopCartRouter.DELETE("/:id", shopcart.Delete)
		ShopCartRouter.POST("", shopcart.New)
		ShopCartRouter.PATCH("/:id", shopcart.Update)
	}
}
