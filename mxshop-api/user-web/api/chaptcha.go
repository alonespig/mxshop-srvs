package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, ans, err := cp.Generate()
	if err != nil {
		zap.L().Error("生成验证码失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成验证码失败"})
		return
	}
	zap.L().Info("生成验证码成功", zap.String("answer", ans))
	c.JSON(http.StatusOK, gin.H{"captchaId": id, "picPath": b64s})
}
