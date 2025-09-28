package handler

import (
	"context"
	"fmt"
	"mxshop/global"
	"mxshop/model"
	"mxshop/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

/*
//商品接口
GoodsList(context.Context, *GoodsFilterRequest) (*GoodsListResponse, error)
//现在用户提交订单有多个商品，你得批量查询商品的信息吧
BatchGetGoods(context.Context, *BatchGoodsIdInfo) (*GoodsListResponse, error)
CreateGoods(context.Context, *CreateGoodsInfo) (*GoodsInfoResponse, error)
DeleteGoods(context.Context, *DeleteGoodsInfo) (*empty.Empty, error)
UpdateGoods(context.Context, *CreateGoodsInfo) (*empty.Empty, error)
GetGoodsDetail(context.Context, *GoodInfoRequest) (*GoodsInfoResponse, error)
*/

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	goodsListResponse := proto.GoodsListResponse{}

	var goods []model.Goods

	localDB := global.DB.Model(&model.Goods{})

	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}

	if req.IsHot {
		localDB = localDB.Where("is_hot = ?", req.IsHot)
	}

	if req.IsNew {
		localDB = localDB.Where("is_new = ?", req.IsNew)
	}

	if req.PriceMin > 0 {
		localDB = localDB.Where("price >= ?", req.PriceMin)
	}

	if req.PriceMax > 0 {
		localDB = localDB.Where("price <= ?", req.PriceMax)
	}

	if req.Brand > 0 {
		localDB = localDB.Where("brand_id = ?", req.Brand)
	}
	var subQuery string
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category where parent_category_id = %d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id = %d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category where id = %d", req.TopCategory)
		}
		localDB = localDB.Where("category_id in (?)", subQuery)
	}

	var count int64
	localDB.Count(&count)

	goodsListResponse.Total = int32(count)

	result := localDB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return &goodsListResponse, nil
}

func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	return nil, status.Errorf(codes.Internal, "implement me")
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	return nil, status.Errorf(codes.Internal, "implement me")
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Internal, "implement me")
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Internal, "implement me")
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	return nil, status.Errorf(codes.Internal, "implement me")
}
