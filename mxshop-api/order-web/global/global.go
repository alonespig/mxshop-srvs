package global

import (
	"mxshop-api/order-web/config"
	"mxshop-api/order-web/proto"
)

var (
	ServerConfig       *config.ServerConfig = &config.ServerConfig{}
	GoodsSrvClient     proto.GoodsClient
	OrderSrvClient     proto.OrderClient
	InventorySrvClient proto.InventoryClient
	NacosConfig        *config.NacosConfig = &config.NacosConfig{}
)
