package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/order-web/api/order"
	"shop-api/order-web/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("orders").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.List)       // 获取订单列表
		OrderRouter.POST("", order.New)       // 新建订单
		OrderRouter.GET("/:id", order.Detail) // 订单详情,可以考虑用订单号
	}
}
