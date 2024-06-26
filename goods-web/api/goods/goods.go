package goods

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"shop-api/goods-web/forms"
	"shop-api/goods-web/global"
	"shop-api/goods-web/proto"
	"shop-api/user-web/api"
)

// List 商品列表
func List(ctx *gin.Context) {
	// 商品的列表
	request := &proto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}
	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	// category
	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("pn", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	// 请求商品服务 goods_srv
	r, err := global.GoodsSrvClient.GoodsList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("[List] 查询商品列表失败", "msg", err.Error())
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": r.Total,
		//"data":  r.Data,
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
			"desc_images": value.DescImages,
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

	// 返回数据
	ctx.JSON(http.StatusOK, reMap)

}

// New 新增商品
func New(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandlerValidatorError(ctx, err)
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
		zap.S().Errorw("[New] 新增商品失败", "msg", err.Error())
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	// todo: 商品库存，和分布式事务一起做
	ctx.JSON(http.StatusOK, rsp)

}

// Detail 商品详情
func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: int32(goodsId)})
	if err != nil {
		zap.S().Errorw("[Detail] 查询商品详情失败", "msg", err.Error())
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	//TODO 去库存服务查询库存
	goodsInfo := map[string]interface{}{
		"id":          rsp.Id,
		"name":        rsp.Name,
		"goods_brief": rsp.GoodsBrief,
		"desc":        rsp.GoodsDesc,
		"ship_free":   rsp.ShipFree,
		"images":      rsp.Images,
		"desc_images": rsp.DescImages,
		"front_image": rsp.GoodsFrontImage,
		"shop_price":  rsp.ShopPrice,
		"category": map[string]interface{}{
			"id":   rsp.Category.Id,
			"name": rsp.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   rsp.Brand.Id,
			"name": rsp.Brand.Name,
			"logo": rsp.Brand.Logo,
		},
		"is_hot":  rsp.IsHot,
		"is_new":  rsp.IsNew,
		"on_sale": rsp.OnSale,
	}
	ctx.JSON(http.StatusOK, goodsInfo)
}

// Delete 删除商品
func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: int32(goodsId)})
	if err != nil {
		zap.S().Errorw("[Delete] 删除商品失败", "msg", err.Error())
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	// ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// Stock 商品库存
func Stock(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// TODO 调用库存服务查询库存
	return
}

func UpdateStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandlerValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(goodsId),
		IsHot:  *goodsStatusForm.IsNew,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandlerValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(goodsId),
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
		zap.S().Errorw("[Update] 更新商品失败", "msg", err.Error())
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}
