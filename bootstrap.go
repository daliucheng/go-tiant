package golib

import (
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/go-tiant/base"
	"github.com/tiant-developer/go-tiant/env"
	"github.com/tiant-developer/go-tiant/flow"
	"github.com/tiant-developer/go-tiant/middleware"
	"github.com/tiant-developer/go-tiant/redis"
	"gorm.io/gorm"
)

type BootstrapConf struct {
	AccessLog      middleware.LoggerConfig `yaml:"accessLog"`
	HandleRecovery gin.RecoveryFunc        `yaml:"handleRecovery"`
	AppName        string                  `yaml:"appName"`
	FlowConfig     FlowConfig              `yaml:"flowConfig"`
	Redis          *redis.Redis
}

type FlowConfig struct {
	Enable  bool
	Db      *gorm.DB
	OutErrs map[int]string
	ErrMap  map[int]int
}

// 全局注册一下，是否使用flow模式开发
func Bootstraps(engine *gin.Engine, conf BootstrapConf) {
	// 环境默认release
	gin.SetMode(env.RunMode)
	// 通用runtime指标采集接口
	base.RegistryMetrics(engine, conf.AppName)
	// access中间键
	engine.Use(middleware.AccessLog(conf.AccessLog))
	// 异常Recovery
	engine.Use(middleware.Recovery(conf.HandleRecovery))
	// api接口监控
	engine.Use(middleware.PromMiddleware(conf.AppName))
	// mvc框架（可选）
	if conf.FlowConfig.Enable {
		flow.SetDefaultDBClient(conf.FlowConfig.Db)
		flow.InitOutErrors(conf.FlowConfig.OutErrs, conf.FlowConfig.ErrMap)
		flow.InitCacheConfig(conf.AppName)
	}
}
