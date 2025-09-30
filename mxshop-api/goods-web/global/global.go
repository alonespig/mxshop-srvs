package global

import (
	"mxshop-api/goods-web/config"
	"mxshop-api/goods-web/proto"
)

var (
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	GoodsSrvClient proto.GoodsClient
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
)
