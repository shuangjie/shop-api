package global

import (
	ut "github.com/go-playground/universal-translator"

	"shop-api/userop-web/config"
	"shop-api/userop-web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}

	GoodsSrvClient proto.GoodsClient
	AddressClient  proto.AddressClient
	MessageClient  proto.MessageClient
	UserFavClient  proto.UserFavClient
)
