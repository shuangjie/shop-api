package main

import (
	"fmt"

	"go.uber.org/zap"

	"shop-api/user-web/initialize"
)

func main() {
	port := 8021
	// 初始化
	initialize.InitLogger()
	router := initialize.Routers()

	zap.S().Infof("user-web服务启动中..., 端口: %d", port)

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("user-web服务启动失败", err.Error())
	}
}
