package global

import (
	"mxshop-api/userop-web/config"
	"mxshop-api/userop-web/proto"
)

var (
	ServerConfig     *config.ServerConfig = &config.ServerConfig{}
	GoodsSrvClient   proto.GoodsClient
	UserFavSrvClient proto.UserFavClient
	MessageSrvClient proto.MessageClient
	AddressSrvClient proto.AddressClient
	NacosConfig      *config.NacosConfig = &config.NacosConfig{}
)
