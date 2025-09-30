package main

import (
	"flag"
	"fmt"
	"net"

	"mxshop/global"
	"mxshop/handler"
	"mxshop/initialize"
	"mxshop/proto"
	"mxshop/utils"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	flag.Parse()

	if *Port == 0 {
		var err error
		*Port, err = utils.GetFreePort()
		if err != nil {
			panic(err)
		}
	}

	zap.L().Debug("ServerConfig", zap.String("ip", *IP), zap.Int("port", *Port), zap.String("name", global.ServerConfig.Name))
	zap.L().Debug("MysqlInfo", zap.Any("mysqlInfo", global.ServerConfig.MysqlInfo))
	zap.L().Debug("ConsulInfo", zap.Any("consulInfo", global.ServerConfig.ConsulInfo))
	// Name       string       `mapstructure:"name" json:"name"`
	// Tags       []string     `mapstructure:"tags" json:"tags"`
	// Host       string       `mapstructure:"host" json:"host"`
	// MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	// ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}

	//注册健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", "172.27.49.67", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)

	if err != nil {
		panic(err)
	}
}
