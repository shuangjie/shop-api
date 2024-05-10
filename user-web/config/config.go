package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"ket"`
}

type AliyunSmsConfig struct {
	RegionID        string `mapstructure:"region_id" json:"region_id"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret"`
	SignName        string `mapstructure:"sign_name" json:"sign_name"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name          string          `mapstructure:"name" json:"name"`
	Port          int             `mapstructure:"port" json:"port"`
	UserSrvInfo   UserSrvConfig   `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo       JWTConfig       `mapstructure:"jwt" json:"jwt"`
	AliyunSmsInfo AliyunSmsConfig `mapstructure:"aliyun_sms" json:"aliyun_sms"`
	RedisInfo     RedisConfig     `mapstructure:"redis" json:"redis"`
	ConsulInfo    ConsulConfig    `mapstructure:"consul" json:"consul"`
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
