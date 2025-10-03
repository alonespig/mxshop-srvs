package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RemoveToStruct(fileds map[string]string) map[string]string {
	rep := make(map[string]string)
	for filed, err := range fileds {
		rep[filed[strings.Index(filed, ".")+1:]] = err
	}
	return rep
}

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

// func HandleValidatorError(c *gin.Context, err error) {
// 	errs, ok := err.(validator.ValidationErrors)
// 	if !ok {
// 		c.JSON(http.StatusOK, gin.H{
// 			"msg": err.Error(),
// 		})
// 	}
// 	c.JSON(http.StatusBadRequest, gin.H{
// 		"error": RemoveToStruct(errs.Translate(global.Trans)),
// 	})
// 	return
// }
