package handler

import (
	"context"
	"log"
	"mxshop/global"
	"mxshop/model"
	"mxshop/proto"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct {
	proto.UnimplementedUserServerServer
}

// type UserServerServer interface {
// 	GetUserList(context.Context, *PageInfo) (*UserListResponse, error)
// 	GetUserByMobile(context.Context, *MobileRequest) (*UserInfoResponse, error)
// 	GetUserById(context.Context, *IdRequest) (*UserInfoResponse, error)
// 	CreateUser(context.Context, *CreateUserInfo) (*UserInfoResponse, error)
// 	UpdateUser(context.Context, *UpdateUserInfo) (*empty.Empty, error)
// 	CheckPassWord(context.Context, *CheckPasswordInfo) (*CheckResponse, error)
// 	mustEmbedUnimplementedUserServerServer()
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
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return &userInfoRsp
}

// GetUserList 获取用户列表
func (u *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}

	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		rsp.Data = append(rsp.Data, ModelToResponse(user))
	}

	return rsp, nil
}

// GetUserByMobile 根据手机号获取用户信息
func (u *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}

// GetUserById 根据ID获取用户信息
func (u *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}

// CreateUser 创建用户
func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected > 0 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName
	// 生成哈希（自动加盐）
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, "创建用户失败")
	}
	return ModelToResponse(user), nil
}

// UpdateUser 更新用户
func (u *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Gender = req.Gender
	user.Birthday = &birthDay

	result = global.DB.Save(&user)

	if result.Error != nil {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}

	return &empty.Empty{}, nil
}

// CheckPassWord 检查密码
func (u *UserServer) CheckPassWord(ctx context.Context, req *proto.CheckPasswordInfo) (*proto.CheckResponse, error) {
	// // 校验密码
	// First argument must be the hashed password, second is the plain text
	log.Println("req.EncryptedPassword", req.EncryptedPassword)
	log.Println("req.Password", req.Password)
	err := bcrypt.CompareHashAndPassword([]byte(req.EncryptedPassword), []byte(req.Password))
	log.Println("err", err)
	return &proto.CheckResponse{
		Success: err == nil,
	}, nil
}

func (u *UserServer) mustEmbedUnimplementedUserServerServer() {
	// panic("unimplemented")
}
