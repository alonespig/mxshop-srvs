package main

import (
	"mxshop-api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {

	//1. 初始化logger
	initialize.InitLogger()
	zap.L().Info("启动服务")

	//2. 初始化路由
	Router := initialize.Routers()

	//3. 启动服务
	if err := Router.Run(":8081"); err != nil {
		zap.L().Error("启动服务失败", zap.Error(err))
		return
	}
}
