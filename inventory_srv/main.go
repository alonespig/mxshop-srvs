package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop/global"
	"mxshop/initialize"
	"mxshop/proto"
	"mxshop/utils"
	"mxshop/utils/register/consul"

	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
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

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &proto.UnimplementedInventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}

	//注册健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	consulInfo := global.ServerConfig.ConsulInfo
	registerClient := consul.NewRegistryClient(consulInfo.Host, consulInfo.Port)
	Id, err := uuid.NewV4()
	if err != nil {
		zap.L().Fatal("生成服务ID失败", zap.Error(err))
	}
	serviceId := fmt.Sprintf("%s", Id)
	serverConfig := global.ServerConfig

	err = registerClient.Register(serverConfig.Host, *Port, serverConfig.Name, serverConfig.Tags, serviceId)
	if err != nil {
		zap.L().Fatal("注册服务失败", zap.Error(err))
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = registerClient.DeRegister(serviceId)
	if err != nil {
		zap.L().Fatal("注销服务失败", zap.Error(err))
	}
	zap.L().Info("服务注销成功")
	zap.L().Info("服务退出")
}
