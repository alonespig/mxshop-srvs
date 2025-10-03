package order

import (
	"mxshop-api/order-web/api"
	"mxshop-api/order-web/forms"
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
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBind(&orderForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(ctx, &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Address: orderForm.Address,
		Mobile:  orderForm.Mobile,
		Name:    orderForm.Name,
		Post:    orderForm.Post,
	})

	if err != nil {
		zap.S().Error("创建订单失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//TODO 返回支付宝到的url
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式错误",
		})
		return
	}

	request := proto.OrderRequest{
		Id: int32(i),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.OrderSrvClient.OrderDetail(ctx, &request)
	if err != nil {
		zap.S().Error("查询订单详情失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"id":       rsp.OrderInfo.Id,
		"status":   rsp.OrderInfo.Status,
		"user":     rsp.OrderInfo.UserId,
		"post":     rsp.OrderInfo.Post,
		"total":    rsp.OrderInfo.Total,
		"address":  rsp.OrderInfo.Address,
		"name":     rsp.OrderInfo.Name,
		"mobile":   rsp.OrderInfo.Mobile,
		"pay_type": rsp.OrderInfo.PayType,
		"order_sn": rsp.OrderInfo.OrderSn,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		goodsList = append(goodsList, map[string]interface{}{
			"id":    item.Id,
			"name":  item.GoodsName,
			"price": item.GoodsPrice,
			"image": item.GoodsImage,
			"nums":  item.Nums,
		})
	}

	reMap["goods"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}
