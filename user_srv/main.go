package main

import (
	"flag"
	"fmt"
	"mxshop/handler"
	"mxshop/initialize"
	"mxshop/proto"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	flag.Parse()

	zap.L().Info("ip:", zap.String("ip", *IP), zap.Int("port", *Port))

	server := grpc.NewServer()
	proto.RegisterUserServerServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
