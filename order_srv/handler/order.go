package handler

import (
	"context"
	"fmt"
	"math/rand"
	"mxshop/global"
	"mxshop/model"
	"mxshop/proto"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver"

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

func (s *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	if result := global.DB.Where("goods = ? and user = ?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车商品不存在")
	}
	return &empty.Empty{}, nil
}

// 更新购物车记录，更新数量和选中状态
func (s *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	var shopCart model.ShoppingCart

	if result := global.DB.Where("goods = ? and user = ?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
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

func generateOrderSn(userId int32) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d%d%d%d%d%d%d%d", userId, time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), rand.Intn(90)+10)
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	//新建订单 访问商品微服务获取商品信息
	//库存的扣减 访问库存微服务
	//从购物车中取到选中的微服务
	// 从购物车中删除已购买的记录
	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "没有选中结算的商品")
	}

	goodsNumMap := make(map[int32]int32)
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumMap[shopCart.Goods] = shopCart.Nums
	}

	//跨服务调用商品 -gin
	goods, err := global.GoodsSrvClient.BatchGetGoods(ctx, &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		return nil, status.Error(codes.Internal, "批量查询商品信息失败")
	}

	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo

	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumMap[good.Id],
		})
	}

	//跨服务调用库存微服务进行库存扣减
	_, err = global.InventorySrvClient.Sell(ctx, &proto.SellInfo{GoodsInfo: goodsInvInfo})
	if err != nil {
		return nil, status.Error(codes.Internal, "库存扣减失败")
	}

	tx := global.DB.Begin()
	//生成订单表
	order := model.OrderInfo{
		OrderSn:      generateOrderSn(req.UserId),
		OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "生成订单失败")
	}

	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}

	//批量插入 orderGoods
	if err := tx.CreateInBatches(orderGoods, 100).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "生成订单失败")
	}

	if err := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "生成订单失败")
	}

	tx.Commit()

	return &proto.OrderInfoResponse{
		Id:      order.ID,
		OrderSn: order.OrderSn,
		Total:   order.OrderMount,
	}, nil
}

func (s *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo

	result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	resp := proto.OrderInfoDetailResponse{
		OrderInfo: &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
		},
	}

	var goods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: req.Id}).Find(&goods); result.Error != nil {
		return nil, result.Error
	}

	for _, good := range goods {
		resp.Goods = append(resp.Goods, &proto.OrderItemResponse{
			GoodsId:    good.Goods,
			GoodsName:  good.GoodsName,
			GoodsPrice: good.GoodsPrice,
			Nums:       good.Nums,
		})
	}
	return &resp, nil
}

func (s *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo

	var total int64
	result := global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := proto.OrderListResponse{}
	rsp.Total = int32(total)

	result = global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &rsp, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*empty.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.Id, req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &empty.Empty{}, nil
}
