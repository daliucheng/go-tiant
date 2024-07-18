package zlog

import (
	"bytes"
	"fmt"
	"github.com/tiant-developer/go-tiant/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// util key
const (
	ContextKeyRequestID = "requestId"
	ContextKeyNoLog     = "_no_log"
	ContextKeyUri       = "_uri"
	zapLoggerAddr       = "_zap_addr"
	sugaredLoggerAddr   = "_sugared_addr"
	customerFieldKey    = "__customerFields"
)

func GetRequestID(ctx *gin.Context) string {
	if ctx == nil {
		return genRequestID()
	}

	// 从ctx中获取
	if r := ctx.GetString(ContextKeyRequestID); r != "" {
		return r
	}
	requestID := genRequestID()
	ctx.Set(ContextKeyRequestID, requestID)
	return requestID
}

var generator = utils.NewRand(time.Now().UnixNano())

func genRequestID() string {
	// 生成 uint64的随机数, 并转换成16进制表示方式
	number := uint64(generator.Int63())
	traceID := fmt.Sprintf("%016x", number)

	var buffer bytes.Buffer
	buffer.WriteString(traceID)
	return buffer.String()
}

func GetRequestUri(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.GetString(ContextKeyUri)
}

// a new method for customer notice
func AddField(c *gin.Context, field ...Field) {
	customerFields := GetCustomerFields(c)
	if customerFields == nil {
		customerFields = field
	} else {
		customerFields = append(customerFields, field...)
	}

	c.Set(customerFieldKey, customerFields)
}

// 获得所有用户自定义的Field
func GetCustomerFields(c *gin.Context) (customerFields []Field) {
	if v, exist := c.Get(customerFieldKey); exist {
		customerFields, _ = v.([]Field)
	}
	return customerFields
}

func SetNoLogFlag(ctx *gin.Context) {
	ctx.Set(ContextKeyNoLog, true)
}

func SetLogFlag(ctx *gin.Context) {
	ctx.Set(ContextKeyNoLog, false)
}

func noLog(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}
	flag, ok := ctx.Get(ContextKeyNoLog)
	if ok && flag == true {
		return true
	}
	return false
}
