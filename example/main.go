package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/golib"
	"github.com/tiant-developer/golib/example/conf"
	"github.com/tiant-developer/golib/example/helpers"
	"github.com/tiant-developer/golib/example/router"
	"github.com/tiant-developer/golib/middleware"
)

func main() {
	defer helpers.Clear()
	// 0.启动前置
	helpers.PreInit()
	// 1.初始化配置和资源
	helpers.Init()
	httpEngine := buildEngine()
	// 3.初始化http服务路由
	router.Http(httpEngine)
	// 服务启动
	fmt.Fprintf(gin.DefaultWriter, "Server start on :%v \n", conf.GlobalC.Port)
	if err := httpEngine.Run(fmt.Sprintf(":%v", conf.GlobalC.Port)); err != nil {
		panic(err.Error())
	}
}

func buildEngine() *gin.Engine {
	engine := gin.New()
	// 2.加载封装的组件

	golib.Bootstraps(engine, golib.BootstrapConf{
		AccessLog: middleware.LoggerConfig{
			SkipCookie:     true,
			SkipPaths:      nil,
			MaxReqBodyLen:  0,
			MaxRespBodyLen: 1024,
			Skip:           nil,
		},
		AppName: conf.GlobalC.AppName,
		Db:      helpers.MysqlClient,
		Redis:   helpers.RedisClient,
	})
	return engine
}
