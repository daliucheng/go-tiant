package test

import (
	"github.com/tiant-developer/golib"
	"github.com/tiant-developer/golib/env"
	"github.com/tiant-developer/golib/example/conf"
	"github.com/tiant-developer/golib/example/helpers"

	"net/http/httptest"
	"path"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
)

var once = sync.Once{}
var Ctx *gin.Context

// Init 基础资源初始化
func Init() {
	once.Do(func() {
		engine := gin.New()
		dir := getSourcePath(0)
		env.SetAppName("testing")
		env.SetRootPath(dir + "/..")
		helpers.PreInit()
		helpers.Init()
		Ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
		golib.Bootstraps(engine, golib.BootstrapConf{
			AppName: conf.GlobalC.AppName,
			Db:      helpers.MysqlClient,
			Redis:   helpers.RedisClient,
		})

	})
}

func getSourcePath(skip int) string {
	_, filename, _, _ := runtime.Caller(skip)
	return path.Dir(filename)
}
