package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"shop-api/order-web/api"
	"shop-api/order-web/forms"
	"shop-api/order-web/global"
	"shop-api/order-web/models"
	"shop-api/order-web/proto"
)

// List 订单列表
func List(ctx *gin.Context) {
	// 1. 获取用户的信息 for jwt
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := proto.OrderFilterRequest{}
	// *：如果是管理员，返回所有订单；如果是用户，返回用户的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	pages := ctx.DefaultQuery("pn", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	// 2. 调用订单服务
	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("[List] 获取 [订单列表] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 3. 返回结果
	reMap := gin.H{
		"total": rsp.Total,
	}
	list := make([]interface{}, 0)
	for _, order := range rsp.Data {
		// 重命名字段
		tmpMap := map[string]interface{}{
			"id":       order.Id,
			"user":     order.UserId,
			"order_sn": order.OrderSn,
			"pay_type": order.PayType,
			"status":   order.Status,
			"total":    order.Total,
			"address":  order.Address,
			"name":     order.Name,
			"mobile":   order.Mobile,
			"post":     order.Post,
			"add_time": order.AddTime,
		}
		list = append(list, tmpMap)
	}
	reMap["data"] = list
	ctx.JSON(http.StatusOK, reMap)
}

// Detail 订单详情
func Detail(ctx *gin.Context) {
	// 1. 获取用户的信息 for jwt
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	// *：如果是管理员，不校验订单的用户id；如果是用户，校验订单的用户id
	request := proto.OrderRequest{
		Id: int32(orderId),
	}
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	// 2. 调用订单服务
	rsp, err := global.OrderSrvClient.OrderDetail(ctx, &request) // mark: 这里 ctx 和 context.Background() 是一样的
	if err != nil {
		zap.S().Errorw("[Detail] 获取 [订单详情] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 3. 拼接返回结果
	reMap := gin.H{
		"id":       rsp.OrderInfo.Id,
		"user":     rsp.OrderInfo.UserId,
		"order_sn": rsp.OrderInfo.OrderSn,
		"pay_type": rsp.OrderInfo.PayType,
		"status":   rsp.OrderInfo.Status,
		"total":    rsp.OrderInfo.Total,
		"address":  rsp.OrderInfo.Address,
		"name":     rsp.OrderInfo.Name,
		"mobile":   rsp.OrderInfo.Mobile,
		"post":     rsp.OrderInfo.Post,
	}

	list := make([]interface{}, 0)
	for _, goods := range rsp.Goods {
		tmpMap := gin.H{
			"id":    goods.GoodsId,
			"name":  goods.GoodsName,
			"price": goods.GoodsPrice,
			"image": goods.GoodsImage,
			"nums":  goods.Nums,
		}
		list = append(list, tmpMap)
	}
	reMap["goods"] = list
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	// 1. 表单验证
	form := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 2. 新建订单
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Address: form.Address,
		Name:    form.Name,
		Mobile:  form.Mobile,
		Post:    form.Post,
	})
	if err != nil {
		zap.S().Errorf("[New] 创建 [订单] 失败: %v", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回结果；todo: 应该返回支付 url
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})

}
