package golib

import (
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/golib/layer"
	"github.com/tiant-developer/golib/middleware"
	"github.com/tiant-developer/golib/redis"
	"gorm.io/gorm"
)

type BootstrapConf struct {
	AccessLog      middleware.LoggerConfig `yaml:"accessLog"`
	HandleRecovery gin.RecoveryFunc        `yaml:"handleRecovery"`
	AppName        string                  `yaml:"appName"`
	Db             *gorm.DB
	Redis          *redis.Redis
}

// 全局注册一下，是否使用flow模式开发
func Bootstraps(engine *gin.Engine, conf BootstrapConf) {
	// 通用runtime指标采集接口
	middleware.RegistryMetrics(engine, conf.AppName)
	// 全局中间键 access日志
	engine.Use(middleware.AccessLog(conf.AccessLog))
	// 异常Recovery
	engine.Use(middleware.Recovery(conf.HandleRecovery))
	// 框架db默认初始化
	layer.SetDefaultDBClient(conf.Db)
	// 框架redis默认初始化
	layer.SetDefaultRedisClient(conf.Redis, conf.AppName)
}
