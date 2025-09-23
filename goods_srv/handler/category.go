package handler

import (
	"context"
	"encoding/json"
	"mxshop/global"
	"mxshop/model"
	"mxshop/proto"

	"github.com/golang/protobuf/ptypes/empty"
)

// 商品分类
func (g *GoodsServer) GetAllCategorysList(ctx context.Context, req *empty.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	data, _ := json.Marshal(categorys)

	return &proto.CategoryListResponse{
		JsonData: string(data),
	}, nil
}

// 获取子分类
func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	return nil, nil
}
func (g *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	return nil, nil
}
func (g *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*empty.Empty, error) {
	return nil, nil
}
func (g *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*empty.Empty, error) {
	return nil, nil
}
