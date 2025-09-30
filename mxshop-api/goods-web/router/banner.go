package router

import (
	"mxshop-api/goods-web/api/banners"
	"mxshop-api/goods-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("", banners.List)                                                          // 轮播图列表页
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.AdminAuth(), banners.Delete) // 删除轮播图
		BannerRouter.POST("", middlewares.JWTAuth(), middlewares.AdminAuth(), banners.New)          //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.AdminAuth(), banners.Update)    //修改轮播图信息
	}
}
