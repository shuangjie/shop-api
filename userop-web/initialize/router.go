package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shop-api/userop-web/middlewares"
	"shop-api/userop-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/up/v1")
	router.InitAddressRouter(ApiGroup)
	router.InitMessageRouter(ApiGroup)
	router.InitUserFavRouter(ApiGroup)

	return Router
}
