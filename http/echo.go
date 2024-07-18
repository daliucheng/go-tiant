package http

import (
	"encoding/json"
	"git.atomecho.cn/atomecho/golib/zlog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EchoXml(ctx *gin.Context, data []byte) {
	ctx.Header("Request_Id", zlog.GetRequestID(ctx))
	ctx.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(data)
}

func EchoXmlError(ctx *gin.Context, data []byte) {
	ctx.Header("Request_Id", zlog.GetRequestID(ctx))
	ctx.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	ctx.Writer.Write(data)
}

func EchoJson(ctx *gin.Context, data []byte) {
	ctx.Header("Request_Id", zlog.GetRequestID(ctx))
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(data)
}

func EchoJsonObj(ctx *gin.Context, data interface{}) {
	ctx.Header("Request_Id", zlog.GetRequestID(ctx))
	dataJson, _ := json.Marshal(data)
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(dataJson)
}

func EchoJsonError(ctx *gin.Context, data []byte) {
	ctx.Header("Request_Id", zlog.GetRequestID(ctx))
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	ctx.Writer.Write(data)
}
