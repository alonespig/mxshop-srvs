package api

import (
	"context"
	"fmt"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 将grpc的code转换成http的状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数错误"})
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "未授权"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "其他错误"})
			}
			return
		}
	}
}

func GetUserList(c *gin.Context) {
	zap.L().Debug("获取用户列表")
	ip := "127.0.0.1"
	port := 50051
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[GetUserList] 连接用户服务失败", zap.Error(err))
		return
	}
	defer userConn.Close()
	//调用接口
	userSrvClient := proto.NewUserServerClient(userConn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil {
		zap.L().Error("[GetUserList] 调用用户服务失败", zap.Error(err))
		HandleGrpcErrorToHttp(err, c)
		return
	}
	result := make([]response.UserResponse, 0)

	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Mobile:   value.Mobile,
			Gender:   value.Gender,
			BirthDay: time.Unix(int64(value.BirthDay), 0).Format("2006-01-02"),
		}
		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)
}
