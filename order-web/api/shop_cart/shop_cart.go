package shop_cart

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"shop-api/order-web/api"
	"shop-api/order-web/global"
	"shop-api/order-web/proto"
)

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
	return
}

func New(ctx *gin.Context) {

}

func Detail(ctx *gin.Context) {

}

func Update(ctx *gin.Context) {

}

func Delete(ctx *gin.Context) {

}
