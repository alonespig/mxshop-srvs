package initialize

import (
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	//设置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("v1")

	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
