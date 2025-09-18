package handler

import (
	"context"
	"mxshop/global"
	"mxshop/model"
	"mxshop/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct {
}

//type UserServerClient interface {
// 	GetUserList(ctx context.Context, in *PageInfo, opts ...grpc.CallOption) (*UserListResponse, error)
// 	GetUserByMobile(ctx context.Context, in *MobileRequest, opts ...grpc.CallOption) (*UserInfoResponse, error)
// 	GetUserById(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*UserInfoResponse, error)
// 	CreateUser(ctx context.Context, in *CreateUserInfo, opts ...grpc.CallOption) (*UserInfoResponse, error)
// 	UpdateUser(ctx context.Context, in *UpdateUserInfo, opts ...grpc.CallOption) (*empty.Empty, error)
// 	CheckPassWord(ctx context.Context, in *CheckPasswordInfo, opts ...grpc.CallOption) (*CheckResponse, error)
// }

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(user model.User) *proto.UserInfoResponse {
	//在grpc的message中字段有默认值，不能随便赋值nil进去，容易出错
	userInfoRsp := proto.UserInfoResponse{
		Id:       int32(user.ID),
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return &userInfoRsp
}

// GetUserList 获取用户列表
func (u *UserServer) GetUserList(ctx context.Context, in *proto.PageInfo) (*proto.UserListResponse, error) {
	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}

	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(in.Pn), int(in.PSize))).Find(&users)

	for _, user := range users {
		rsp.Data = append(rsp.Data, ModelToResponse(user))
	}

	return rsp, nil
}

// GetUserByMobile 根据手机号获取用户信息
func (u *UserServer) GetUserByMobile(ctx context.Context, in *proto.MobileRequest, opts ...grpc.CallOption) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: in.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}

// GetUserById 根据ID获取用户信息
func (u *UserServer) GetUserById(ctx context.Context, in *proto.IdRequest, opts ...grpc.CallOption) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, in.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}
