package flow

import (
	"github.com/gin-gonic/gin"
	"github.com/tiant-developer/go-tiant/errors"
	"github.com/tiant-developer/go-tiant/flow/layer"
	"github.com/tiant-developer/go-tiant/utils"
	"gorm.io/gorm"
	"reflect"
)

func slave(src interface{}) interface{} {
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

func UseController(ctl layer.IController) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		controller := slave(ctl).(layer.IController)
		controller.SetContext(ctx)
		controller.SetEntity(controller)
		controller.OnCreate(nil)
		controllerName := utils.ReflectType(controller).Name()
		// get request dto
		req, bindType := controller.BindReq()
		if req == nil {
			controller.LogErrorf("Controller %s has no DtoRequest struct", controllerName)
			controller.RenderJsonFail(errors.ErrorSystemError)
			return
		}
		var err error
		if bindType == nil {
			err = ctx.ShouldBind(req)
		} else {
			err = ctx.ShouldBindWith(req, bindType)
		}
		if err != nil {
			controller.LogError("Controller %s param bind error, err:%s", controllerName, err.Error())
			controller.RenderJsonFail(errors.ErrorParamInvalid)
			return
		}
		// action execute
		data, err := controller.Action()
		if err != nil {
			controller.LogError("Controller %s call action logic error, err:%s", controllerName, err.Error())
			controller.RenderJsonFail(err)
			return
		}
		// 支持自定义渲染
		if controller.ShouldRender() {
			controller.RenderJsonSucc(data)
		}
	}
}

func UseFlow(ctx *gin.Context, newFlow layer.ILayer) interface{} {
	return UseFlowParam(ctx, newFlow, nil)
}

func UseFlowParam(ctx *gin.Context, newFlow layer.ILayer, param layer.IFlowParam) interface{} {
	newFlow.SetContext(ctx)
	newFlow.SetEntity(newFlow)
	newFlow.OnCreate(param)
	return newFlow
}

// *****以下为需要初始化执行的方法*******
// 初始化默认数据库
func SetDefaultDBClient(db *gorm.DB) {
	layer.SetDefaultDBClient(db)
}

// 初始化带名称数据库（可支持多数据库）
func SetNamedDBClient(namedDbs map[string]*gorm.DB) {
	layer.SetNamedDBClient(namedDbs)
}

// 初始化错误码，准出错误码和业务&准出错误码映射
func InitOutErrors(outErrs map[int]string, errMap map[int]int) {
	errors.InitOutErrors(outErrs, errMap)
}

// 初始化缓存配置，业务线和模块名称
func InitCacheConfig(appName string) {
	layer.InitCacheConfig(appName)
}

// 设置是否关闭自动读主库，默认是开启的，这是总开关
func SetCloseAutoReadDBMaster(close bool) {
	layer.SetCloseAutoReadMaster(close)
}
