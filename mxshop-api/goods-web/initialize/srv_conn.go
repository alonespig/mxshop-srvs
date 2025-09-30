package initialize

import (
	"fmt"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // 注册 resolver
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal("[InitSrvConn] 连接用户服务失败", zap.Error(err))
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
}
