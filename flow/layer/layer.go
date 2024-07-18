package layer

import (
	"github.com/tiant-developer/go-tiant/utils"
	"github.com/tiant-developer/go-tiant/zlog"

	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
	"time"
)

type IFlowParam interface {
	Validate() error
}

type FlowParam struct {
}

func (param *FlowParam) Validate() error {
	return nil
}

type ILayer interface {
	GetContext() *gin.Context
	SetContext(*gin.Context)
	Create(newFlow ILayer) interface{}
	CreateWithParam(newFlow ILayer, param IFlowParam) interface{}
	Assign(targets ...interface{})
	AssignWithParam(dst interface{}, param IFlowParam)
	CopyWithContext(ctx *gin.Context) interface{}
	OnCreate(param IFlowParam)
	SetEntity(entity ILayer)
	GetEntity() ILayer
	StartTimer(timerKey string)
	StopTimer(timerKey string) int
	SetReadDbMaster(isReadMaster bool)

	LogDebug(args ...interface{})
	LogDebugf(format string, args ...interface{})
	LogInfo(args ...interface{})
	LogInfof(format string, args ...interface{})
	LogWarn(args ...interface{})
	LogWarnf(format string, args ...interface{})
	LogError(args ...interface{})
	LogErrorf(format string, args ...interface{})
}

type Layer struct {
	ctx    *gin.Context
	entity ILayer
	Param  IFlowParam

	m     sync.Mutex
	timer map[string]int64
}

func (entity *Layer) SetContext(ctx *gin.Context) {
	entity.ctx = ctx
}

func (entity *Layer) GetContext() *gin.Context {
	return entity.ctx
}

func (entity *Layer) SetEntity(flow ILayer) {
	entity.entity = flow
}

func (entity *Layer) GetEntity() ILayer {
	return entity.entity
}

func (entity *Layer) OnCreate(param IFlowParam) {
	entity.Param = param
}

func (entity *Layer) Create(newFlow ILayer) interface{} {
	return entity.CreateWithParam(newFlow, nil)
}

func (entity *Layer) CreateWithParam(newFlow ILayer, param IFlowParam) interface{} {
	newFlow.SetContext(entity.ctx)
	newFlow.SetEntity(newFlow)
	newFlow.OnCreate(param)
	return newFlow
}

func Create(ctx *gin.Context, newFlow ILayer) interface{} {
	return createWithParam(ctx, newFlow, nil)
}

func createWithParam(ctx *gin.Context, newFlow ILayer, param IFlowParam) interface{} {
	newFlow.SetContext(ctx)
	newFlow.SetEntity(newFlow)
	newFlow.OnCreate(param)
	return newFlow
}

func (entity *Layer) Assign(targets ...interface{}) {
	// 遍历，根据target的类型new出对象，并赋值到target指针
	for _, dst := range targets {
		entity.AssignWithParam(dst, nil)
	}
}

func (entity *Layer) AssignWithParam(dst interface{}, param IFlowParam) {
	pDst := reflect.ValueOf(dst)
	if pDst.Kind() == reflect.Ptr {
		pDst = pDst.Elem()
	}
	if pDst.Kind() == reflect.Ptr {
		t := pDst.Type().Elem()
		v := reflect.New(t).Elem().Addr().Interface().(ILayer)
		flow := entity.CreateWithParam(v, param)
		pDst.Set(reflect.ValueOf(flow))
	}
}

func (entity *Layer) CopyWithContext(ctx *gin.Context) interface{} {
	v := utils.NewObject(entity.entity).(ILayer)
	e := entity.CreateWithParam(v, entity.Param).(ILayer)
	if ctx != nil {
		e.SetContext(ctx)
	}
	return e
}

// 标记是否需要读主库
func (entity *Layer) SetReadDbMaster(isReadMaster bool) {
	entity.ctx.Set("__isReadDbMaster__", isReadMaster)
}

func (entity *Layer) GetReadDbMaster() bool {
	if v, exist := entity.ctx.Get("__isReadDbMaster__"); exist {
		if is, ok := v.(bool); ok {
			return is
		}
	}
	return false
}

// 标记是否强制读主库，用于覆盖全局开关
func (entity *Layer) ForceReadDbMaster(forceReadMaster bool) {
	entity.ctx.Set("__forceReadDbMaster__", forceReadMaster)
}

func (entity *Layer) IsForceReadDbMaster() bool {
	if v, exist := entity.ctx.Get("__forceReadDbMaster__"); exist {
		if is, ok := v.(bool); ok {
			return is
		}
	}
	return false
}

func (entity *Layer) StartTimer(timerKey string) {
	entity.m.Lock()
	defer entity.m.Unlock()
	if entity.timer == nil {
		entity.timer = make(map[string]int64)
	}
	entity.timer[timerKey] = time.Now().UnixNano()
}

func (entity *Layer) StopTimer(timerKey string) int {
	entity.m.Lock()
	defer entity.m.Unlock()
	if entity.timer == nil {
		return 0
	}
	if v, ok := entity.timer[timerKey]; ok {
		now := time.Now().UnixNano()
		pass := int((now - v) / int64(time.Millisecond)) //ms
		zlog.AddField(entity.GetContext(), zlog.Int(timerKey, pass))
		delete(entity.timer, timerKey)
		return pass
	}
	return 0
}

// 日志打印方法
func (entity *Layer) LogDebug(args ...interface{}) {
	zlog.Debug(entity.ctx, args...)
}

func (entity *Layer) LogDebugf(format string, args ...interface{}) {
	zlog.Debugf(entity.ctx, format, args...)
}

func (entity *Layer) LogInfo(args ...interface{}) {
	zlog.Info(entity.ctx, args...)
}

func (entity *Layer) LogInfof(format string, args ...interface{}) {
	zlog.Infof(entity.ctx, format, args...)
}

func (entity *Layer) LogWarn(args ...interface{}) {
	zlog.Warn(entity.ctx, args...)
}

func (entity *Layer) LogWarnf(format string, args ...interface{}) {
	zlog.Warnf(entity.ctx, format, args...)
}

func (entity *Layer) LogError(args ...interface{}) {
	zlog.Error(entity.ctx, args...)
}

func (entity *Layer) LogErrorf(format string, args ...interface{}) {
	zlog.Errorf(entity.ctx, format, args...)
}
