package initialize

import (
	"mxshop-api/goods-web/middlewares"
	"mxshop-api/goods-web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	//设置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/g/v1")

	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)

	return Router
}
