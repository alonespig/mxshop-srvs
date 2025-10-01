package main

import (
	"context"
	"fmt"
	"mxshop/proto"

	"google.golang.org/grpc"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	brandClient = proto.NewGoodsClient(conn)
}

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func main() {
	Init()
	defer conn.Close()
	TestGetBrandList()
}
