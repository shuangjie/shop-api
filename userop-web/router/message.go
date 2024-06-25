package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/userop-web/api/message"

	"shop-api/userop-web/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List) // 获取留言列表
		MessageRouter.POST("", message.New) // 新增留言
	}
}
