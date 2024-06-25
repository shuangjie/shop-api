package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"shop-api/userop-web/global"
	"shop-api/userop-web/initialize"
	"shop-api/userop-web/utils"
	"shop-api/userop-web/utils/register/consul"
	validator2 "shop-api/userop-web/validator"
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

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", validator2.ValidateMobile)
		v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := uuid.NewV4().String()
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("注册服务失败", err.Error())
	}

	zap.S().Debugf("userop-web服务启动中..., 端口: %d", global.ServerConfig.Port)
	go func() {
		if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("userop-web服务启动失败", err.Error())
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销服务失败", err.Error())
	} else {
		zap.S().Info("userop-web服务注销成功")
	}

}
