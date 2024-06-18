package global

import (
	ut "github.com/go-playground/universal-translator"

	"shop-api/order-web/config"
	"shop-api/order-web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}

	OrderSrvClient proto.OrderClient
	GoodsSrvClient proto.GoodsClient
)
