package initialize

import (
	"mxshop-api/order-web/middlewares"
	"mxshop-api/order-web/router"
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

	ApiGroup := Router.Group("/o/v1")

	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)

	return Router
}
