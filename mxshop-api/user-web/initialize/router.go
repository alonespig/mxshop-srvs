package initialize

import (
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("v1")

	router.InitUserRouter(ApiGroup)

	return Router
}
