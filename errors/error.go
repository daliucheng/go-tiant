package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type Error struct {
	ErrNo  int
	ErrMsg string
}

func NewError(code int, message string) *Error {
	return &Error{
		ErrNo:  code,
		ErrMsg: message,
	}
}

func (err Error) Error() string {
	return err.ErrMsg
}

func (err Error) Sprintf(v ...interface{}) Error {
	err.ErrMsg = fmt.Sprintf(err.ErrMsg, v...)
	return err
}

func (err Error) Equal(e error) bool {
	switch errors.Cause(e).(type) {
	case Error:
		return err.ErrNo == errors.Cause(e).(Error).ErrNo
	default:
		return false
	}
}

func (err Error) WrapPrint(core error, message string, user ...interface{}) error {
	if core == nil {
		return nil
	}
	err.SetErrPrintfMsg(core)
	return errors.Wrap(err, message)
}

func (err Error) WrapPrintf(core error, format string, message ...interface{}) error {
	if core == nil {
		return nil
	}
	err.SetErrPrintfMsg(core)
	return errors.Wrap(err, fmt.Sprintf(format, message...))
}

func (err Error) Wrap(core error) error {
	if core == nil {
		return nil
	}

	msg := err.ErrMsg
	err.ErrMsg = core.Error()
	return errors.Wrap(err, msg)
}

func (err *Error) SetErrPrintfMsg(v ...interface{}) {
	err.ErrMsg = fmt.Sprintf(err.ErrMsg, v...)
}

// 标准准出错误码定义
const (
	// 通用错误码
	PARAM_ERROR     = 1   //参数错误
	SYSTEM_ERROR    = 2   //服务内部错误
	USER_NOT_LOGIN  = 3   //用户未登录
	INVALID_REQUEST = 6   //无效请求
	DEFAULT_ERROR   = 100 //默认错误，未准出的错误码，会修改为此错误码
	CUSTOM_ERROR    = 101 //自定义错误，无固定错误文案

	// 业务错误码，规范：6位数字长度，前2位模块，后6位业务自定义
)

// 标准准出错误码文案定义
var ErrMsg = map[int]string{
	// 通用错误文案
	PARAM_ERROR:     "请求参数错误",
	SYSTEM_ERROR:    "服务异常，请稍后重试",
	USER_NOT_LOGIN:  "用户Session已失效，请重新登录",
	INVALID_REQUEST: "请求无效，请稍后再试",
	DEFAULT_ERROR:   "服务开小差了，请稍后再试",
}

// *****以下是通用准出错误码的简便定义***********
// 正常
var ErrorSuccess = Error{
	ErrNo:  0,
	ErrMsg: "success",
}

// 参数错误
var ErrorParamInvalid = Error{
	ErrNo:  PARAM_ERROR,
	ErrMsg: ErrMsg[PARAM_ERROR],
}

// 系统异常
var ErrorSystemError = Error{
	ErrNo:  SYSTEM_ERROR,
	ErrMsg: ErrMsg[SYSTEM_ERROR],
}

// 用户未登录
var ErrorUserNotLogin = Error{
	ErrNo:  USER_NOT_LOGIN,
	ErrMsg: ErrMsg[USER_NOT_LOGIN],
}

// 无效请求
var ErrorInvalidRequest = Error{
	ErrNo:  INVALID_REQUEST,
	ErrMsg: ErrMsg[INVALID_REQUEST],
}

// 默认错误
var ErrorDefault = Error{
	ErrNo:  DEFAULT_ERROR,
	ErrMsg: ErrMsg[DEFAULT_ERROR],
}

// 自定义错误
// 使用方式：ErrorCustomError.Sprintf(v)
var ErrorCustomError = Error{
	ErrNo:  CUSTOM_ERROR,
	ErrMsg: "%s",
}

// 错误码转换映射表，用以支持业务模块错误码和标准准出错误码映射
var ErrNoMap = map[int]int{}

func FormatOutputError(e error) error {
	e = errors.Cause(e)
	err, ok := e.(Error)
	if !ok {
		return Error{
			ErrNo:  DEFAULT_ERROR,
			ErrMsg: ErrMsg[DEFAULT_ERROR],
		}
	}
	if err.ErrNo == CUSTOM_ERROR {
		return Error{
			ErrNo:  CUSTOM_ERROR,
			ErrMsg: err.ErrMsg,
		}
	}
	// 错误码转换
	errNo, exist := ErrNoMap[err.ErrNo]
	if !exist {
		errNo = err.ErrNo
	}
	// 标准化输出
	if errMsg, exist := ErrMsg[errNo]; exist {
		return Error{
			ErrNo:  errNo,
			ErrMsg: errMsg,
		}
	}
	return Error{
		ErrNo:  DEFAULT_ERROR,
		ErrMsg: ErrMsg[DEFAULT_ERROR],
	}
}

// 校验业务准出错误码是否符合规范（6位长度）
func checkOutErrorFormat(errNo int) bool {
	if errNo < 1000000 || errNo > 999999 {
		return false
	}
	return true
}

// 初始化业务准出错误码和映射关系。注意：准出的错误码需要符合规范，同时更新到wiki
// outErrs-业务定义的标准准出错误码
// errMap-业务使用的错误码和标准准出错误码映射关系
func InitOutErrors(outErrs map[int]string, errMap map[int]int) {
	// 补充标准准出错误码
	for errNo, errMsg := range outErrs {
		if !checkOutErrorFormat(errNo) {
			continue
		}
		ErrMsg[errNo] = errMsg
	}
	// 补充错误码转换映射
	for errNoIn, errNoOut := range errMap {
		// 检查errNoOut是否标准准出错误码，非准出错误码忽略
		if _, exist := ErrMsg[errNoOut]; !exist {
			continue
		}
		ErrNoMap[errNoIn] = errNoOut
	}
}
