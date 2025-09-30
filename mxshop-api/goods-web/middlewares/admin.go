package middlewares

import (
	"mxshop-api/goods-web/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		zap.L().Debug("AdminAuth")
		claims, _ := c.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId == 2 {
			c.JSON(http.StatusForbidden, gin.H{"msg": "权限不足"})
			c.Abort()
			return
		}
		c.Next()
	}
}
