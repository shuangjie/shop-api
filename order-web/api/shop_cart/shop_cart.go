package shop_cart

import (
	"context"
	"net/http"
	"shop-api/order-web/forms"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"shop-api/order-web/api"
	"shop-api/order-web/global"
	"shop-api/order-web/proto"
)

// List 获取购物车列表
func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 [购物车列表] 失败")
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}
	// 如果购物车为空，直接返回
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// 查询购物车中商品的详细信息
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 [商品信息] 失败")
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	// 组装数据
	data := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, goods := range goodsRsp.Data {
			if item.GoodsId == goods.Id {
				tmp := map[string]interface{}{
					"id":          item.Id,
					"goods_id":    item.GoodsId,
					"goods_name":  goods.Name,
					"goods_image": goods.GoodsFrontImage,
					"goods_price": goods.ShopPrice,
					"nums":        item.Nums,
					"checked":     item.Checked,
				}
				data = append(data, tmp)
			}
		}
	}

	reMap := gin.H{
		"total": rsp.Total,
		"data":  data,
	}

	ctx.JSON(http.StatusOK, reMap)
}

// New 添加购物车
func New(ctx *gin.Context) {
	// 1. 参数校验
	itemForm := forms.ShopCartForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandlerValidatorError(ctx, err)
		return
	}

	// 2. 查询商品信息
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[New] 查询 [商品信息] 失败")
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	// 3. 查询库存信息（库存是否充足）
	invRsp, err := global.InventSrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[New] 查询 [库存] 失败")
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "库存不足",
		})
		return
	}

	// 4. 添加购物车
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("[New] 添加 [购物车] 失败")
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Detail(ctx *gin.Context) {

}

func Update(ctx *gin.Context) {

}

// Delete 删除购物车
func Delete(ctx *gin.Context) {
	// 获取商品 id
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	// 删除购物车
	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: int32(i),
		UserId:  userId.(int32),
	})
	if err != nil {
		zap.S().Errorw("[Delete] 删除 [购物车] 失败")
		api.HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
