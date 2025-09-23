package handler

import (
	"mxshop/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// 商品接口
// func (g *GoodsServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method GoodsList not implemented")
// }

// func (g *GoodsServer) BatchGetGoods(context.Context, *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method BatchGetGoods not implemented")
// }

// func (g *GoodsServer) CreateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method CreateGoods not implemented")
// }

// func (g *GoodsServer) DeleteGoods(context.Context, *proto.DeleteGoodsInfo) (*empty.Empty, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method DeleteGoods not implemented")
// }

// func (g *GoodsServer) UpdateGoods(context.Context, *proto.CreateGoodsInfo) (*empty.Empty, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method UpdateGoods not implemented")
// }

// func (g *GoodsServer) GetGoodsDetail(context.Context, *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method GetGoodsDetail not implemented")
// }
