package goods

import (
	"context"
	"mxshop-api/goods-web/forms"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func removeToStruct(fileds map[string]string) map[string]string {
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

func List(c *gin.Context) {
	priceMin := c.DefaultQuery("priceMin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)

	priceMax := c.DefaultQuery("priceMax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)

	isHot := c.DefaultQuery("isHot", "0")

	isNew := c.DefaultQuery("in", "0")

	isTab := c.DefaultQuery("isTab", "0")

	categoryId := c.DefaultQuery("categoryId", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)

	page := c.DefaultQuery("page", "1")
	pageInt, _ := strconv.Atoi(page)

	pageSize := c.DefaultQuery("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)

	keywords := c.DefaultQuery("q", "")

	brandId := c.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)

	request := &proto.GoodsFilterRequest{
		PriceMin:    int32(priceMinInt),
		PriceMax:    int32(priceMaxInt),
		IsHot:       isHot == "1",
		IsNew:       isNew == "1",
		IsTab:       isTab == "1",
		TopCategory: int32(categoryIdInt),
		KeyWords:    keywords,
		Pages:       int32(pageInt),
		PagePerNums: int32(pageSizeInt),
		Brand:       int32(brandIdInt),
	}

	//请求商品的服务
	r, err := global.GoodsSrvClient.GoodsList(context.Background(), request)
	if err != nil {
		zap.L().Error("[List]请求商品的服务失败", zap.Error(err))
		HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_image":  value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}

	reMap["data"] = goodsList

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		zap.L().Error("[New] 绑定商品表单失败", zap.Error(err))
		HandleGrpcErrorToHttp(err, c)
		return
	}

	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	//TODO 商品的库存
	c.JSON(http.StatusOK, rsp)
}

func Detail(c *gin.Context) {
	goodsId := c.Param("id")
	goodsIdInt, err := strconv.Atoi(goodsId)

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	goodsClient := global.GoodsSrvClient
	r, err := goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: int32(goodsIdInt),
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	rspMap := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"ctegory": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}
	c.JSON(http.StatusOK, rspMap)
}

func Delete(c *gin.Context) {
	goodsId := c.Param("id")
	goodsIdInt, err := strconv.Atoi(goodsId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	goodsClient := global.GoodsSrvClient
	_, err = goodsClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: int32(goodsIdInt),
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func Stocks(c *gin.Context) {
	goodsId := c.Param("id")
	goodsIdInt, err := strconv.Atoi(goodsId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	zap.L().Info("[Stocks] 获取商品库存", zap.Int("goodsId", goodsIdInt))

	// TODO 商品的库存
	return
}

func UpdateStatus(c *gin.Context) {
	goodsForm := forms.GoodsStatusForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		zap.L().Error("[UpdateStatus] 绑定商品状态表单失败", zap.Error(err))
		HandleGrpcErrorToHttp(err, c)
		return
	}

	goodsId := c.Param("id")
	goodsIdInt, err := strconv.Atoi(goodsId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	goodsClient := global.GoodsSrvClient
	_, err = goodsClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(goodsIdInt),
		IsHot:  *goodsForm.IsHot,
		OnSale: *goodsForm.OnSale,
		IsNew:  *goodsForm.IsNew,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		zap.L().Error("[Update] 绑定商品表单失败", zap.Error(err))
		HandleGrpcErrorToHttp(err, c)
		return
	}

	goodsId := c.Param("id")
	goodsIdInt, err := strconv.Atoi(goodsId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	goodsClient := global.GoodsSrvClient
	_, err = goodsClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(goodsIdInt),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
