package layer

import "fmt"

const (
	EXPIRE_TIME_1_SECOND  = 1
	EXPIRE_TIME_5_SECOND  = 5
	EXPIRE_TIME_30_SECOND = 30
	EXPIRE_TIME_1_MINUTE  = 60
	EXPIRE_TIME_5_MINUTE  = 300
	EXPIRE_TIME_15_MINUTE = 900
	EXPIRE_TIME_30_MINUTE = 1800
	EXPIRE_TIME_1_HOUR    = 3600
	EXPIRE_TIME_2_HOUR    = 7200
	EXPIRE_TIME_6_HOUR    = 21600
	EXPIRE_TIME_12_HOUR   = 43200
	EXPIRE_TIME_1_DAY     = 86400
	EXPIRE_TIME_3_DAY     = 259200
	EXPIRE_TIME_1_WEEK    = 604800
)

var ModuleName string //模块

func InitCacheConfig(moduleName string) {
	ModuleName = moduleName
}

type IRedis interface {
	ILayer
	RedisFunc()
}

type Redis struct {
	Layer
}

func (entity *Redis) RedisFunc() {
	fmt.Print("this is redis func\n")
}

func (entity *Redis) FormatCacheKey(format string, args ...interface{}) string {
	prefix := getPrefix()
	return fmt.Sprintf(prefix+format, args...)
}

func getPrefix() string {
	prefix := ""
	// 模块名默认module，很容易冲突
	if ModuleName == "" {
		prefix += "module:"
	} else {
		prefix += ModuleName + ":"
	}
	return prefix
}
