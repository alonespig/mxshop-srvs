package initialize

import (
	"fmt"
	"mxshop/global"
	"mxshop/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// 初始化第三方微服务的client
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

	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal("[InitSrvConn] 连接库存服务失败", zap.Error(err))
	}
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}
