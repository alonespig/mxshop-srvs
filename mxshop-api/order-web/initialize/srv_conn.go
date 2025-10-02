package initialize

import (
	"fmt"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"

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
		zap.L().Fatal("[InitSrvConn] 连接商品服务失败", zap.Error(err))
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal("[InitSrvConn] 连接订单服务失败", zap.Error(err))
	}

	global.OrderSrvClient = proto.NewOrderClient(orderConn)
}
