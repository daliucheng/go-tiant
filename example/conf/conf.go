package conf

import (
	"github.com/tiant-developer/golib/api"
	"github.com/tiant-developer/golib/env"
	"github.com/tiant-developer/golib/orm"
	"github.com/tiant-developer/golib/redis"
	"github.com/tiant-developer/golib/zlog"
)

var (
	// 全局配置变量
	GlobalC SResource
)

type SResource struct {
	Port    int            `yaml:"port"`
	SDKPort int            `yaml:"sdkPort"`
	AppName string         `yaml:"appName"`
	Log     zlog.LogConfig `yaml:"log"`
	Mysql   map[string]orm.MysqlConf
	Redis   map[string]redis.RedisConf
	Api     map[string]*api.ApiClient // 调用三方后台
	Oidc    SOidc                     `yaml:"oidc"`
	Kafka   SKafKa                    `yaml:"kafka"`
}

type SOidc struct {
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	IssuerURL    string `yaml:"issuer_url"`
}

type SKafKa struct {
	ProxyUrl      string `yaml:"proxyUrl"`
	ProxyUserName string `yaml:"proxyUserName"`
	ProxyPassword string `yaml:"proxyPassword"`
	GroupId       string `yaml:"groupId"`
	Topic         string `yaml:"topic"`
	StartSync     bool   `yaml:"startSync"`
}

func InitConf() {
	// 从环境变量加载资源类配置
	env.LoadConf("resource.yaml", "mount", &GlobalC)
}
