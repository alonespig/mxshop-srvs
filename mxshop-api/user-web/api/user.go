package api

import (
	"context"
	"fmt"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[GetUserList] 连接用户服务失败", zap.Error(err))
		return
	}
	defer userConn.Close()
	//调用接口
	userSrvClient := proto.NewUserServerClient(userConn)

	pn := c.DefaultQuery("pn", "1")
	psize := c.DefaultQuery("psize", "10")

	pnInt, _ := strconv.ParseInt(pn, 10, 32)
	psizeInt, _ := strconv.ParseInt(psize, 10, 32)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(psizeInt),
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

// PassWordLogin 密码登录
func PassWordLogin(c *gin.Context) {
	passWordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBindJSON(&passWordLoginForm); err != nil {
		zap.L().Error("[PassWordLogin] 绑定参数失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("[PassWordLogin] 连接用户服务失败", zap.Error(err))
		return
	}
	defer userConn.Close()
	userSrvClient := proto.NewUserServerClient(userConn)

	rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	})

	if err != nil {
		zap.L().Error("[PassWordLogin] 调用用户服务失败", zap.Error(err))
		HandleGrpcErrorToHttp(err, c)
		return
	}

	passRep, err := userSrvClient.CheckPassWord(context.Background(), &proto.CheckPasswordInfo{
		Password:          passWordLoginForm.PassWord,
		EncryptedPassword: rsp.Password,
	})
	if err != nil {
		zap.L().Error("[PassWordLogin] 调用用户服务失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "登录失败，密码错误"})
		return
	}
	if passRep.Success {
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(rsp.Id),
			NickName:    rsp.NickName,
			AuthorityId: uint(rsp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + 60*60*24*7,
				Issuer:    "mxshop-api",
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":        rsp.Id,
			"nick_name": rsp.NickName,
			"token":     token,
			"expire":    (time.Now().Unix() + 60*60*24*7) * 1000,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "登录失败，密码错误"})
	}
}
