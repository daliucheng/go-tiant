package middleware

import (
	"bytes"
	"fmt"
	"github.com/tiant-developer/go-tiant/utils"
	"github.com/tiant-developer/go-tiant/zlog"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	_defaultPrintRequestLen  = 10240
	_defaultPrintResponseLen = 10240
)

type customRespWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w customRespWriter) WriteString(s string) (int, error) {
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w customRespWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

// access日志打印
type LoggerConfig struct {
	// SkipPaths is a url path array which logs are not written.
	SkipPaths []string `yaml:"skipPaths"`

	// request body 最大长度展示，0表示采用默认的10240，-1表示不打印
	MaxReqBodyLen int `yaml:"maxReqBodyLen"`
	// response body 最大长度展示，0表示采用默认的10240，-1表示不打印。指定长度的时候需注意，返回的json可能被截断
	MaxRespBodyLen int `yaml:"maxRespBodyLen"`

	// 自定义Skip功能
	Skip func(ctx *gin.Context) bool
}

func AccessLog(conf LoggerConfig) gin.HandlerFunc {
	notLogged := conf.SkipPaths
	var skip map[string]struct{}
	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	maxReqBodyLen := conf.MaxReqBodyLen
	if maxReqBodyLen == 0 {
		maxReqBodyLen = _defaultPrintRequestLen
	}

	maxRespBodyLen := conf.MaxRespBodyLen
	if maxRespBodyLen == 0 {
		maxRespBodyLen = _defaultPrintResponseLen
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// body writer
		blw := &customRespWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 请求参数，涉及到回写，要在处理业务逻辑之前
		reqParam := getReqBody(c, maxReqBodyLen)

		c.Set(zlog.ContextKeyUri, path)
		_ = zlog.GetRequestID(c)

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; ok {
			return
		}

		if conf.Skip != nil && conf.Skip(c) {
			return
		}

		// Stop timer
		end := time.Now()

		response := ""
		if blw.body != nil && maxRespBodyLen != -1 {
			response = blw.body.String()
			if len(response) > maxRespBodyLen {
				response = response[:maxRespBodyLen]
			}
		}

		// 固定notice
		commonFields := []zlog.Field{
			zlog.String("method", c.Request.Method),
			zlog.String("clientIp", utils.GetClientIp(c)),
			zlog.String("cookie", getCookie(c)),
			zlog.String("reqStartTime", utils.GetFormatRequestTime(start)),
			zlog.String("reqEndTime", utils.GetFormatRequestTime(end)),
			zlog.Float64("cost", utils.GetRequestCost(start, end)),
			zlog.String("requestParam", reqParam),
			zlog.Int("responseStatus", c.Writer.Status()),
			zlog.String("response", response),
			zlog.Int("bodySize", c.Writer.Size()),
			zlog.String("reqErr", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		}

		// 新的notice添加方式
		customerFields := zlog.GetCustomerFields(c)
		commonFields = append(commonFields, customerFields...)
		zlog.InfoLogger(c, "notice", commonFields...)
	}
}

// 请求参数
func getReqBody(c *gin.Context, maxReqBodyLen int) (reqBody string) {
	// 不打印参数
	if maxReqBodyLen == -1 {
		return reqBody
	}

	// body中的参数
	if c.Request.Body != nil && c.ContentType() == binding.MIMEMultipartPOSTForm {
		requestBody, err := c.GetRawData()
		if err != nil {
			zlog.WarnLogger(c, "get http request body error: "+err.Error())
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		if _, err := c.MultipartForm(); err != nil {
			zlog.WarnLogger(c, "parse http request form body error: "+err.Error())
		}
		reqBody = c.Request.PostForm.Encode()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	} else if c.Request.Body != nil && c.ContentType() == "application/octet-stream" {

	} else if c.Request.Body != nil {
		requestBody, err := c.GetRawData()
		if err != nil {
			zlog.WarnLogger(c, "get http request body error: "+err.Error())
		}
		reqBody = utils.BytesToString(requestBody)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// 拼接上 url?rawQuery 的参数 todo 为了兼容以前逻辑，感觉参数应该分开写更好?
	if c.Request.URL.RawQuery != "" {
		reqBody += "&" + c.Request.URL.RawQuery
	}

	// 截断参数
	if len(reqBody) > maxReqBodyLen {
		reqBody = reqBody[:maxReqBodyLen]
	}

	return reqBody
}

// 从request body中解析特定字段作为notice key打印
func getReqValueByKey(ctx *gin.Context, k string) string {
	if vs, exist := ctx.Request.Form[k]; exist && len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func getCookie(ctx *gin.Context) string {
	cStr := ""
	for _, c := range ctx.Request.Cookies() {
		cStr += fmt.Sprintf("%s=%s&", c.Name, c.Value)
	}
	return strings.TrimRight(cStr, "&")
}

func AddField(field ...zlog.Field) gin.HandlerFunc {
	return func(c *gin.Context) {
		zlog.AddField(c, field...)
		c.Next()
	}
}
