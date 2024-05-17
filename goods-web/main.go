package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"shop-api/goods-web/global"
	"shop-api/goods-web/initialize"
	"shop-api/goods-web/utils"
	"shop-api/goods-web/utils/register/consul"
)

func main() {
	// 初始化日志
	initialize.InitLogger()
	// 初始化配置
	initialize.InitConfig()
	// 初始化翻译器
	router := initialize.Routers()

	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Panic("初始化翻译器失败", err.Error())
	}

	// 初始化srv连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	// 本地开发环境端口号固定，生产环境则随机获取
	env := viper.GetBool("SHOP_ENV_DEV")
	if !env {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, uuid.NewV4().String())
	if err != nil {
		zap.S().Panic("注册服务失败", err.Error())
	}

	zap.S().Infof("goods-web服务启动中..., 端口: %d", global.ServerConfig.Port)
	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("goods-web服务启动失败", err.Error())
	}
}
