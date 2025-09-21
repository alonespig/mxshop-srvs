package initialize

import (
	"fmt"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	//从注册中心获取用户服务地址
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	zap.L().Debug("Address", zap.Any("cfg.Address", cfg.Address))

	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.L().Error("[GetUserList] 创建consul客户端失败", zap.Error(err))
		return
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		zap.L().Error("[GetUserList] 获取用户服务地址失败", zap.Error(err))
		return
	}

	zap.L().Debug("data", zap.Any("data", data))

	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == "" || userSrvPort == 0 {
		zap.L().Fatal("[InitSrvConn] 获取用户服务地址失败")
		return
	}

	zap.L().Debug("userSrvHost", zap.Any("userSrvHost", userSrvHost))
	zap.L().Debug("userSrvPort", zap.Any("userSrvPort", userSrvPort))

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[InitSrvConn] 连接用户服务失败", zap.Error(err))
		return
	}
	//调用接口
	userSrvClient := proto.NewUserServerClient(userConn)
	global.UserSrvClient = userSrvClient
}
