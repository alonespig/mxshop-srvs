package order

import (
	"mxshop-api/order-web/api"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/models"
	"mxshop-api/order-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	model := claims.(models.CustomClaims)

	request := proto.OrderFilterRequest{}

	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	pages := ctx.DefaultQuery("p", "1")
	pageInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pageInt)

	pageSize := ctx.DefaultQuery("page_size", "10")
	perNumsInt, _ := strconv.Atoi(pageSize)
	request.PagePerNums = int32(perNumsInt)

	rsp, err := global.OrderSrvClient.OrderList(ctx, &request)
	if err != nil {
		zap.S().Error("查询订单失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	repMap := gin.H{
		"total": rsp.Total,
	}

	orderList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{
			"id":       item.Id,
			"status":   item.Status,
			"pay_type": item.PayType,
			"user":     item.Id,
			"post":     item.Post,
			"total":    item.Total,
			"address":  item.Address,
			"name":     item.Name,
			"mobile":   item.Mobile,
			"order_sn": item.OrderSn,
			"add_time": item.AddTime,
		}
		orderList = append(orderList, tmpMap)
	}

	repMap["data"] = orderList

	ctx.JSON(http.StatusOK, repMap)
}

func New(ctx *gin.Context) {

}

func Detail(ctx *gin.Context) {

}
