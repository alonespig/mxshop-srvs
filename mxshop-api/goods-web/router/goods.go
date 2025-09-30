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
		GoodsRouter.POST("/list", middlewares.JWTAuth(), middlewares.AdminAuth(), goods.List)
	}
}
