package main

import (
	"flag"
	"fmt"
	"mxshop/handler"
	"mxshop/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 8080, "端口号")

	flag.Parse()
	fmt.Println("ip:", *IP, "port:", *Port)

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
