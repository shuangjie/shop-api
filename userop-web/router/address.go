package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/userop-web/api/address"

	"shop-api/userop-web/middlewares"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address").Use(middlewares.JWTAuth())
	{
		AddressRouter.GET("", address.List)          // 获取地址列表
		AddressRouter.DELETE("/:id", address.Delete) // 删除地址
		AddressRouter.POST("", address.New)          // 新增地址
		AddressRouter.PUT("/:id", address.Update)    // 更新地址
	}
}
