package config

type SrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"ket"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type AlipayConfig struct {
	AppID           string `mapstructure:"app_id" json:"app_id"`
	PrivateKey      string `mapstructure:"private_key" json:"private_key"`
	AlipayPublicKey string `mapstructure:"alipay_public_key" json:"alipay_public_key"`
	IsProduction    bool   `mapstructure:"is_production" json:"is_production"`
	NotifyURL       string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL       string `mapstructure:"return_url" json:"return_url"`
}

type ServerConfig struct {
	Name             string       `mapstructure:"name" json:"name"`
	Host             string       `mapstructure:"host" json:"host"`
	Port             int          `mapstructure:"port" json:"port"`
	Tags             []string     `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo     SrvConfig    `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo SrvConfig    `mapstructure:"inventory_srv" json:"inventory_srv"`
	OrderSrvInfo     SrvConfig    `mapstructure:"order_srv" json:"order_srv"`
	JWTInfo          JWTConfig    `mapstructure:"jwt" json:"jwt"`
	ConsulInfo       ConsulConfig `mapstructure:"consul" json:"consul"`
	AlipayInfo       AlipayConfig `mapstructure:"alipay" json:"alipay"`
}

type NacosConfig struct {
	Host        string `mapstructure:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespace_id"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Group       string `mapstructure:"group"`
	DataId      string `mapstructure:"data_id"`
}
