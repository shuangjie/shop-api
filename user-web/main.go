package main

import (
	"fmt"

	"go.uber.org/zap"

	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
)

func main() {
	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	router := initialize.Routers()

	zap.S().Infof("user-web服务启动中..., 端口: %d", global.ServerConfig.Port)

	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("user-web服务启动失败", err.Error())
	}
}
