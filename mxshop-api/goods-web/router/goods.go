package router

import (
	"mxshop-api/goods-web/api/goods"
	"mxshop-api/user-web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	zap.L().Info("配置用户相关的url")
	{
		GoodsRouter.GET("/list", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.AdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.AdminAuth(), goods.Delete)
		GoodsRouter.GET("/:id/stock", goods.Stocks)

		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.AdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.AdminAuth(), goods.UpdateStatus)
	}
}
