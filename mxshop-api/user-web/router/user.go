package router

import (
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.L().Info("配置用户相关的url")
	{
		UserRouter.GET("/list", middlewares.JWTAuth(), middlewares.AdminAuth(), api.GetUserList)
		UserRouter.POST("/login", api.PassWordLogin)
		UserRouter.POST("/register", api.Register)
	}
}
