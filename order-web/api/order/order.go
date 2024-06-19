package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"shop-api/order-web/api"
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
		api.HandlerGrpcErrorToHttp(err, ctx)
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

func New(context *gin.Context) {

}

func Detail(context *gin.Context) {

}
