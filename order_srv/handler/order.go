package handler

import (
	"context"
	"mxshop/global"
	"mxshop/model"
	"mxshop/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func (s *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	//获取用户的购物车列表
	var shopCarts []model.ShoppingCart

	result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := proto.CartItemListResponse{}
	rsp.Total = int32(result.RowsAffected)
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}
func (s *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车
	var shopCart model.ShoppingCart

	result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&shopCart)

	if result.RowsAffected > 0 {
		shopCart.Nums += req.Nums
	} else {
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	result = global.DB.Save(&shopCart)
	if result.Error != nil {
		return nil, result.Error
	}

	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}
func (s *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {

	return nil, nil
}
func (s *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.ShoppingCart{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车商品不存在")
	}
	return &empty.Empty{}, nil
}
func (s *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	return nil, nil
}
func (s *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	return nil, nil
}

// 更新购物车记录，更新数量和选中状态
func (s *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	var shopCart model.ShoppingCart

	if result := global.DB.First(&shopCart, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车商品不存在")
	}
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	result := global.DB.Save(&shopCart)
	if result.Error != nil {
		return nil, result.Error
	}

	return &empty.Empty{}, nil
}
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*empty.Empty, error) {
	return nil, nil
}
