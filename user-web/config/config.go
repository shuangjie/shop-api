package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type AliyunSmsConfig struct {
	RegionID        string `mapstructure:"region_id"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	SignName        string `mapstructure:"sign_name"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string          `mapstructure:"name"`
	Port          int             `mapstructure:"port"`
	UserSrvInfo   UserSrvConfig   `mapstructure:"user_srv"`
	JWTInfo       JWTConfig       `mapstructure:"jwt"`
	AliyunSmsInfo AliyunSmsConfig `mapstructure:"aliyun_sms"`
	RedisInfo     RedisConfig     `mapstructure:"redis"`
}
