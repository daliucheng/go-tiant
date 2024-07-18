package layer

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tiant-developer/go-tiant/api"
	"github.com/tiant-developer/go-tiant/errors"
)

// 接口重要级别定义，关系日志打印
type ApiLevel int

const API_LEVEL_CRITICAL ApiLevel = 1
const API_LEVEL_NORMAL ApiLevel = 2
const API_LEVEL_EASY ApiLevel = 3

type Res struct {
	//兼容errNo errno
	ErrNo int `json:"errNo"`
	Errno int `json:"errno"`
	//兼容errMsg、errmsg、errStr、errstr
	ErrMsg interface{}         `json:"errMsg"`
	Errmsg interface{}         `json:"errmsg"`
	ErrStr interface{}         `json:"errStr"`
	Errstr interface{}         `json:"errstr"`
	Data   jsoniter.RawMessage `json:"data"`
}

type ApiRes struct {
	ErrNo  int
	ErrMsg string
	Data   []byte
}

type IApi interface {
	ILayer
	ApiFunc()
	GetEncodeType() string
	GetApiLevel(path string, errNo int) ApiLevel
	ApiGet(path string, requestParam interface{}) (*ApiRes, error)
	ApiPost(path string, requestBody interface{}) (*ApiRes, error)
	ApiGetWithOpts(path string, reqOpts api.HttpRequestOptions) (*ApiRes, error)
	ApiPostWithOpts(path string, reqOpts api.HttpRequestOptions) (*ApiRes, error)
}

type Api struct {
	Layer
	EncodeType string
	Client     *api.ApiClient
}

func (entity *Api) ApiFunc() {
	fmt.Print("this is api func\n")
}

// api请求数据格式，默认json
func (entity *Api) GetEncodeType() string {
	if entity.EncodeType != "" {
		return entity.EncodeType
	}
	return api.EncodeJson
}

// 默认所有接口都是critical，需要打error日志，具体接口和错误码可以重写此方法
func (entity *Api) GetApiLevel(path string, errNo int) ApiLevel {
	return API_LEVEL_CRITICAL
}

func (entity *Api) ApiGet(path string, requestParam interface{}) (*ApiRes, error) {
	reqOpts := api.HttpRequestOptions{
		RequestBody: requestParam,
		Encode:      api.EncodeForm,
	}
	return entity.ApiGetWithOpts(path, reqOpts)
}

func (entity *Api) ApiPost(path string, requestBody interface{}) (*ApiRes, error) {
	api2 := entity.GetEntity().(IApi)
	reqOpts := api.HttpRequestOptions{
		RequestBody: requestBody,
		Encode:      api2.GetEncodeType(),
	}
	return entity.ApiPostWithOpts(path, reqOpts)
}

func (entity *Api) ApiGetWithOpts(path string, reqOpts api.HttpRequestOptions) (*ApiRes, error) {
	//GET请求写死为form
	reqOpts.Encode = api.EncodeForm

	if entity.Client == nil {
		entity.LogErrorf("ApiGetWithOpts failed, api client is needed, path:%s", path)
		return nil, errors.ErrorSystemError
	}
	res, e := entity.Client.HttpGet(entity.GetContext(), path, reqOpts)
	if e != nil {
		entity.LogErrorf("ApiGetWithOpts failed, path:%s, err:%v", path, e)
		return nil, e
	}
	return entity.handel(path, res)
}

func (entity *Api) ApiPostWithOpts(path string, reqOpts api.HttpRequestOptions) (*ApiRes, error) {
	if reqOpts.Encode == "" {
		api := entity.GetEntity().(IApi)
		reqOpts.Encode = api.GetEncodeType()
	}

	if entity.Client == nil {
		entity.LogErrorf("ApiPostWithOpts failed, api client is needed, path:%s", path)
		return nil, errors.ErrorSystemError
	}
	res, e := entity.Client.HttpPost(entity.GetContext(), path, reqOpts)
	if e != nil {
		entity.LogErrorf("ApiPostWithOpts failed, path:%s, err:%v", path, e)
		return nil, e
	}
	return entity.handel(path, res)
}

func (entity *Api) handel(path string, res *api.ApiResult) (*ApiRes, error) {
	httpRes := Res{}
	e := jsoniter.Unmarshal(res.Response, &httpRes)
	if e != nil {
		// 限制一下错误日志打印的长度，2k
		data := res.Response
		if len(data) > 2000 {
			data = data[0:2000]
		}
		// 返回数据json unmarshal失败，打印错误日志
		entity.LogErrorf("http response json unmarshal failed, path:%s, response:%s, err:%v", path, string(data), e)
		return nil, e
	}
	//兼容各式各样的错误码
	errNo, errMsg := compatibleErr(&httpRes)
	if errNo != 0 {
		// errNo != 0，是否打印日志取决于接口级别
		api := entity.GetEntity().(IApi)
		apiLevel := api.GetApiLevel(path, errNo)
		switch apiLevel {
		case API_LEVEL_CRITICAL:
			entity.LogErrorf("rpc call get origin data has error, path:%s, errNo:%d, errMsg:%s", path, errNo, errMsg)
		case API_LEVEL_NORMAL:
			entity.LogWarnf("rpc call get origin data has error, path:%s, errNo:%d, errMsg:%s", path, errNo, errMsg)
		case API_LEVEL_EASY:
		default:
		}
	}
	apiRes := &ApiRes{
		ErrNo:  errNo,
		ErrMsg: errMsg,
		Data:   httpRes.Data,
	}
	return apiRes, nil
}

// Api中封装json unmarshal 2个方法，支持自动打日志功能，业务方调用此方法无需再打错误日志
func (entity *Api) Unmarshal(data []byte, v interface{}) error {
	err := jsoniter.Unmarshal(data, v)
	if err != nil {
		// 限制一下错误日志打印的长度，2k
		if len(data) > 2000 {
			data = data[0:2000]
		}
		entity.LogErrorf("data json unmarshal failed, data:%s, err:%v", string(data), err)
	}
	return err
}

func (entity *Api) UnmarshalFromString(data string, v interface{}) error {
	err := jsoniter.UnmarshalFromString(data, v)
	if err != nil {
		// 限制一下错误日志打印的长度，2k
		if len(data) > 2000 {
			data = data[0:2000]
		}
		entity.LogErrorf("data json unmarshal failed, data:%s, err:%v", data, err)
	}
	return err
}

// 此方法同Unmarshal的区别为兼容了data=[]的情况，但注意out未初始化可能为nil，应用方按需调用两种方法
func (entity *Api) ParseJsonData(data []byte, out interface{}) error {
	// 兼容php的data=[]情况
	if len(data) == 2 && string(data) == "[]" {
		return nil
	}
	err := jsoniter.Unmarshal(data, out)
	if err != nil {
		// 限制一下错误日志打印的长度，2k
		if len(data) > 2000 {
			data = data[0:2000]
		}

		entity.LogErrorf("data json unmarshal failed, data:%s, err:%v", data, err)
	}
	return err
}

func compatibleErr(res *Res) (int, string) {
	errNo := 0
	errMsg := ""
	if res.ErrNo != 0 {
		//errNo
		errNo = res.ErrNo
	} else if res.Errno != 0 {
		//errno
		errNo = res.Errno
	}

	if msg, ok := res.ErrStr.(string); msg != "" && ok {
		//errMsg
		errMsg = msg
	} else if msg, ok := res.ErrMsg.(string); msg != "" && ok {
		//errmsg
		errMsg = msg
	} else if msg, ok := res.Errstr.(string); msg != "" && ok {
		//errStr
		errMsg = msg
	} else if msg, ok := res.Errmsg.(string); msg != "" && ok {
		//errstr
		errMsg = msg
	}

	return errNo, errMsg
}

func (entity *Api) DecodeApiResponse(outPut interface{}, data *ApiRes, err error) error {
	if err != nil {
		entity.LogErrorf("api error, http request failed, err : %s", err.Error())
		return errors.ErrorSystemError
	}

	if data.ErrNo != 0 {
		return errors.Error{
			ErrNo:  data.ErrNo,
			ErrMsg: data.ErrMsg,
		}
	}

	// 解析数据
	if err = entity.Unmarshal(data.Data, outPut); err != nil {
		entity.LogErrorf("api error, api response unmarshal err : %s", err.Error())
		return errors.ErrorSystemError
	}

	return nil
}
