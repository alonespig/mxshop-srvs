package router

import (
	"mxshop-api/order-web/api/order"
	"mxshop-api/order-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order")
	{
		OrderRouter.GET("", middlewares.JWTAuth(), middlewares.AdminAuth(), order.List)
		OrderRouter.POST("", middlewares.JWTAuth(), order.New)
		OrderRouter.GET("/:id", middlewares.JWTAuth(), order.Detail)
	}
}
