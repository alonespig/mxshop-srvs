package shopcart

import (
	"context"
	"mxshop-api/order-web/forms"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"
	"net/http"
	"strconv"

	"mxshop-api/order-web/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId.(int32),
	})
	if err != nil {
		zap.S().Error("查询购物车失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.Id)
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Error("查询商品失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	repMap := gin.H{
		"total": rsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if item.GoodsId == good.Id {
				rmpMap := map[string]interface{}{
					"id":          item.Id,
					"goods_id":    item.GoodsId,
					"goods_name":  good.Name,
					"goods_image": good.GoodsFrontImage,
					"goods_price": good.ShopPrice,
					"nums":        item.Nums,
					"checked":     item.Checked,
				}
				goodsList = append(goodsList, rmpMap)
			}
		}
	}
	repMap["data"] = goodsList

	ctx.JSON(http.StatusOK, repMap)
}

func New(ctx *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := ctx.ShouldBind(&itemForm); err != nil {
		zap.S().Error("绑定失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//为了严谨，添加商品到购物车之前，先检查商品是否存在
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Error("查询商品失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Error("查询商品失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "库存不足",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		Nums:    itemForm.Nums,
		UserId:  int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Error("添加商品到购物车失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error("转换id失败", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式错误",
		})
		return
	}

	userId, _ := ctx.Get("userId")

	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		Id:     int32(i),
		UserId: int32(userId.(uint)),
	})

	if err != nil {
		zap.S().Error("删除商品失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		Id:     int32(i),
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Error("删除商品失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error("转换id失败", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式错误",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		zap.S().Error("绑定失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	userId, _ := ctx.Get("userId")

	request := proto.CartItemRequest{
		UserId:  userId.(int32),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}

	if itemForm.Checked == nil {
		request.Checked = *itemForm.Checked
	}

	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if err != nil {
		zap.S().Error("更新商品失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
