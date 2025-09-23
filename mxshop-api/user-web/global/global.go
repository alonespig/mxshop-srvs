package global

import (
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	UserSrvClient proto.UserServerClient
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
)
