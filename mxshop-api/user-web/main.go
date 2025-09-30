package main

import (
	"fmt"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	"mxshop-api/user-web/utils/register/consul"
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

	//服务注册
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
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.L().Error("启动服务失败", zap.Error(err))
			return
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
