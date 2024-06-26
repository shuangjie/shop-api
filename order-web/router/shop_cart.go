package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/order-web/api/shop_cart"
	"shop-api/order-web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shop_cart").Use(middlewares.JWTAuth())
	{
		ShopCartRouter.GET("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), shop_cart.List) // 获取购物车列表
		ShopCartRouter.POST("", middlewares.JWTAuth(), shop_cart.New)                            // 新建购物车
		ShopCartRouter.PATCH("/:id", middlewares.JWTAuth(), shop_cart.Update)                    // 更新购物车
		ShopCartRouter.DELETE("/:id", middlewares.JWTAuth(), shop_cart.Delete)                   // 删除购物车 这个 id 是商品的id
	}
}
