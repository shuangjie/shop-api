package router

import (
	"github.com/gin-gonic/gin"

	"shop-api/goods-web/api/goods"
	"shop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)                                                                 // 商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)               // 新建商品
		GoodsRouter.GET("/:id", goods.Detail)                                                           // 商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)      // 删除商品
		GoodsRouter.GET("/:id/stock", goods.Stock)                                                      // 商品库存
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         // 更新商品
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus) // 更新商品状态
	}
}
