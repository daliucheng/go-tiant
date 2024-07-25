package layer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/go-tiant/errors"
	http2 "github.com/tiant-developer/go-tiant/http"
	"net/http"
	"reflect"
)

type ControllerErrType string

type IController[T any] interface {
	ILayer
	Action(req T) (any, error)
	ShouldRender() bool
	RenderJsonFail(err error)
	RenderJsonSuccess(data any)
}

type Controller struct {
	Layer
}

func (entity *Controller) Action(any) (any, error) {
	return nil, nil
}

func (entity *Controller) ShouldRender() bool {
	return true
}

func (entity *Controller) RenderJsonFail(err error) {
	http2.RenderJsonFail(entity.GetCtx(), err)
}

func (entity *Controller) RenderJsonSuccess(data any) {
	http2.RenderJsonSucc(entity.GetCtx(), data)
}

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
	rander := http2.DefaultRender{}
	if e, ok := err.(errors.Error); ok {
		rander.Code = e.Code
		rander.Message = e.Message
	} else {
		rander.Code = errors.ErrorSystemError.Code
		rander.Message = errors.ErrorSystemError.Message
	}
	flusher, _ := ctx.Writer.(http.Flusher)
	str, _ := json.Marshal(rander)
	fmt.Fprintf(ctx.Writer, "%s\n", str)
	flusher.Flush()
}
func slave(src any) any {
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { //如果是指针类型
		typ = typ.Elem()               //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem() //创建对象
		return dst.Addr().Interface()  //返回指针
	} else {
		dst := reflect.New(typ).Elem() //创建对象
		return dst.Interface()         //返回值
	}
}

func Use[T any](controller IController[*T]) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		newCTL := slave(controller).(IController[*T])
		var newReq T
		newCTL.SetCtx(ctx)
		newCTL.SetEntity(controller)
		if len(ctx.ContentType()) == 0 && ctx.Request.Method == http.MethodPost { // post默认application/json
			err := ctx.BindJSON(&newReq)
			if err != nil {
				newCTL.LogWarnf("Controller %T param bind error, err:%+v", newCTL, err)
				newCTL.RenderJsonFail(errors.ErrorParamInvalid)
				return
			}
		} else {
			err := ctx.ShouldBind(&newReq)
			if err != nil {
				newCTL.LogWarnf("Controller %T param bind error, err:%+v", newCTL, err)
				newCTL.RenderJsonFail(errors.ErrorParamInvalid)
				return
			}
		}
		// action execute
		data, err := newCTL.Action(&newReq)
		if err != nil {
			newCTL.LogWarnf("Controller %T call action logic error, err:%+v", newCTL, err)
			newCTL.RenderJsonFail(err)
			return
		}
		// 支持自定义渲染
		if newCTL.ShouldRender() {
			newCTL.RenderJsonSuccess(data)
		}
	}
}
