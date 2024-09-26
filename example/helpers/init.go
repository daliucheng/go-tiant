package helpers

import (
	"github.com/tiant-developer/golib/env"
	"github.com/tiant-developer/golib/example/conf"
	"github.com/tiant-developer/golib/zlog"
)

// 基础资源（必须）
func PreInit() {
	// 动态加载nacos配置覆盖默认配置
	conf.InitConf()
	// 用于日志中展示模块的名字手动指定
	env.SetAppName(conf.GlobalC.AppName)
	// 初始化zlog日志
	zlog.InitLog(conf.GlobalC.AppName, conf.GlobalC.Log)
}

func Init() {
	InitMysql()
	InitRedis()

}

func Clear() {
	// 服务结束时的清理工作，对应 Init() 初始化的资源
	zlog.CloseLogger()
	//CloseRedis()
}
