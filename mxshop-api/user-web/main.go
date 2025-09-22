package main

import (
	"fmt"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"

	"github.com/spf13/viper"
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

	initialize.InitSrvConn()

	viper.AutomaticEnv()
	debug := viper.GetBool("MXSHOP_DEBUG")
	// 如果是本地开发环境，端口号固定
	if debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	//4. 启动服务
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.L().Error("启动服务失败", zap.Error(err))
		return
	}
}
