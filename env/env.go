package env

import (
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/go-tiant/utils"
	"os"
	"path/filepath"
)

const DefaultRootPath = "."

var (
	// 本地ip
	LocalIP string
	// 根目录
	rootPath string
	// 是否docker运行
	isDocker bool
	// 项目AppName
	AppName string
	// gin运行模式
	RunMode string
)

func init() {
	LocalIP = utils.GetLocalIp()
	isDocker = false
	// 运行环境
	RunMode = gin.ReleaseMode
	r := os.Getenv("DOCKER_ENV")
	switch r {
	case "online":
		RunMode = gin.ReleaseMode
		isDocker = true
	default:
		RunMode = gin.DebugMode
	}
}

// RootPath 返回应用的根目录
func GetRootPath() string {
	if rootPath != "" {
		return rootPath
	} else {
		return DefaultRootPath
	}
}

// GetConfDirPath 返回配置文件目录绝对地址
func GetConfDirPath() string {
	return filepath.Join(GetRootPath(), "conf")
}

// LogRootPath 返回log目录的绝对地址
func GetLogDirPath() string {
	return filepath.Join(GetRootPath(), "log")
}

// 判断项目运行平台
func IsDockerPlatform() bool {
	return isDocker
}

// 手动指定SetAppName
func SetAppName(appName string) {
	AppName = appName
}

func GetAppName() string {
	return AppName
}

// SetRootPath 设置应用的根目录
func SetRootPath(r string) {
	if !isDocker {
		rootPath = r
	}
}
