package main

import (
	"fmt"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/initialize"
	"mxshop-api/order-web/utils"
	"mxshop-api/order-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
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
	consulInfo := global.ServerConfig.ConsulInfo
	registerClient := consul.NewRegistryClient(consulInfo.Host, consulInfo.Port)
	Id, err := uuid.NewV4()
	if err != nil {
		zap.L().Fatal("生成服务ID失败", zap.Error(err))
	}
	serviceId := fmt.Sprintf("%s", Id)
	serverConfig := global.ServerConfig

	err = registerClient.Register(serverConfig.Host, serverConfig.Port, serverConfig.Name, serverConfig.Tags, serviceId)
	if err != nil {
		zap.L().Fatal("注册服务失败", zap.Error(err))
	}
	//4. 启动服务

	zap.S().Debugf("启动服务，端口：%d", global.ServerConfig.Port)

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.L().Fatal("启动服务失败", zap.Error(err))
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//注销服务
	err = registerClient.DeRegister(serviceId)
	if err != nil {
		zap.L().Fatal("注销服务失败", zap.Error(err))
	} else {
		zap.L().Info("注销服务成功")
	}
	zap.L().Info("服务退出")
}
