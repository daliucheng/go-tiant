package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/go-tiant/base"
	"github.com/tiant-developer/go-tiant/errors"
	"net/http"
)

func EchoXml(ctx *gin.Context, data []byte) {
	ctx.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(data)
}

func EchoXmlError(ctx *gin.Context, data []byte) {
	ctx.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	ctx.Writer.Write(data)
}

func EchoJson(ctx *gin.Context, data []byte) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(data)
}

func EchoJsonObj(ctx *gin.Context, data interface{}) {
	dataJson, _ := json.Marshal(data)
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(dataJson)
}

func EchoJsonError(ctx *gin.Context, data []byte) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	ctx.Writer.Write(data)
}

func EchoStreamError(ctx *gin.Context, err error) {
	rander := base.DefaultRender{}
	if e, ok := err.(errors.Error); ok {
		rander.ErrNo = e.ErrNo
		rander.ErrMsg = e.ErrMsg
	} else {
		rander.ErrNo = errors.ErrorSystemError.ErrNo
		rander.ErrMsg = errors.ErrorSystemError.ErrMsg
	}
	flusher, _ := ctx.Writer.(http.Flusher)
	str, _ := json.Marshal(rander)
	fmt.Fprintf(ctx.Writer, "%s\n", str)
	flusher.Flush()
}
