package main

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
	validator2 "shop-api/user-web/validator"
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

	zap.S().Infof("user-web服务启动中..., 端口: %d", global.ServerConfig.Port)

	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("user-web服务启动失败", err.Error())
	}
}
