package router

import (
	"mxshop-api/order-web/api/order"
	"mxshop-api/order-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.List)
		OrderRouter.POST("", order.New)
		OrderRouter.GET("/:id", order.Detail)
	}
}
