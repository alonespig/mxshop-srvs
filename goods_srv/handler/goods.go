package handler

import (
	"mxshop/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

/*
type GoodsServer interface {
	//商品接口
	GoodsList(context.Context, *GoodsFilterRequest) (*GoodsListResponse, error)
	//现在用户提交订单有多个商品，你得批量查询商品的信息吧
	BatchGetGoods(context.Context, *BatchGoodsIdInfo) (*GoodsListResponse, error)
	CreateGoods(context.Context, *CreateGoodsInfo) (*GoodsInfoResponse, error)
	DeleteGoods(context.Context, *DeleteGoodsInfo) (*empty.Empty, error)
	UpdateGoods(context.Context, *CreateGoodsInfo) (*empty.Empty, error)
	GetGoodsDetail(context.Context, *GoodInfoRequest) (*GoodsInfoResponse, error)
	//商品分类
	GetAllCategorysList(context.Context, *empty.Empty) (*CategoryListResponse, error)
	//获取子分类
	GetSubCategory(context.Context, *CategoryListRequest) (*SubCategoryListResponse, error)
	CreateCategory(context.Context, *CategoryInfoRequest) (*CategoryInfoResponse, error)
	DeleteCategory(context.Context, *DeleteCategoryRequest) (*empty.Empty, error)
	UpdateCategory(context.Context, *CategoryInfoRequest) (*empty.Empty, error)
	//品牌和轮播图
	BrandList(context.Context, *BrandFilterRequest) (*BrandListResponse, error)
	CreateBrand(context.Context, *BrandRequest) (*BrandInfoResponse, error)
	DeleteBrand(context.Context, *BrandRequest) (*empty.Empty, error)
	UpdateBrand(context.Context, *BrandRequest) (*empty.Empty, error)
	//轮播图
	BannerList(context.Context, *empty.Empty) (*BannerListResponse, error)
	CreateBanner(context.Context, *BannerRequest) (*BannerResponse, error)
	DeleteBanner(context.Context, *BannerRequest) (*empty.Empty, error)
	UpdateBanner(context.Context, *BannerRequest) (*empty.Empty, error)
	//品牌分类
	CategoryBrandList(context.Context, *CategoryBrandFilterRequest) (*CategoryBrandListResponse, error)
	//通过category获取brands
	GetCategoryBrandList(context.Context, *CategoryInfoRequest) (*BrandListResponse, error)
	CreateCategoryBrand(context.Context, *CategoryBrandRequest) (*CategoryBrandResponse, error)
	DeleteCategoryBrand(context.Context, *CategoryBrandRequest) (*empty.Empty, error)
	UpdateCategoryBrand(context.Context, *CategoryBrandRequest) (*empty.Empty, error)
	mustEmbedUnimplementedGoodsServer()
}
*/
