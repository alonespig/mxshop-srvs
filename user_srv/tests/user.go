package main

import (
	"context"
	"fmt"
	"mxshop/proto"

	"google.golang.org/grpc"
)

var userClient proto.UserServerClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserServerClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    2,
		PSize: 10,
	})
	if err != nil {
		panic(err)
	}

	for _, user := range rsp.Data {
		fmt.Println(user)
		rsp, err := userClient.CheckPassWord(context.Background(), &proto.CheckPasswordInfo{
			Password:          "123456",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp)
	}
}

func main() {
	Init()
	defer conn.Close()
	TestGetUserList()
}
