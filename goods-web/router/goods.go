package router

import (
	"github.com/gin-gonic/gin"

	"shop-api/goods-web/api/goods"
	"shop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)                                                   // 商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New) // 新建商品
		GoodsRouter.GET("/:id", goods.Detail)                                             // 商品详情
	}
}
