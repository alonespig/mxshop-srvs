package main

import (
	"fmt"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {

	//1. 初始化logger
	initialize.InitLogger()
	zap.L().Info("启动服务")

	//3. 初始化配置
	initialize.InitConfig()

	//2. 初始化路由
	Router := initialize.Routers()

	//4. 启动服务
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.L().Error("启动服务失败", zap.Error(err))
		return
	}
}
