package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"github.com/r3labs/sse/v2"
	"github.com/tiant-developer/go-tiant/errors"
	"net/http"
)

// 定义 SSE 事件
type MessageEvent struct {
	Id    string
	Event string
	Data  string
}

// 实现 SSE 事件的 String() 方法
func (e MessageEvent) String() string {
	return fmt.Sprintf("id:%s\n"+
		"event:%s\n"+
		"data:%s\n\n", e.Id, e.Event, e.Data)
}

// 流式输出报错
func EchoStreamError(ctx *gin.Context, err error) {
	rander := DefaultRender{}
	if e, ok := err.(errors.Error); ok {
		rander.ErrNo = e.ErrNo
		rander.ErrMsg = e.ErrMsg
	} else {
		rander.ErrNo = errors.ErrorSystemError.ErrNo
		rander.ErrMsg = errors.ErrorSystemError.ErrMsg
	}
	flusher, _ := ctx.Writer.(http.Flusher)
	str, _ := json.Marshal(rander)
	msg := MessageEvent{
		Id:    "",
		Event: "error",
		Data:  string(str),
	}
	fmt.Fprintf(ctx.Writer, "%s", msg.String())
	flusher.Flush()
}

func EchoStream(ctx *gin.Context, id, event, str string) {
	flusher, _ := ctx.Writer.(http.Flusher)
	msg := MessageEvent{
		Id:    id,
		Event: event,
		Data:  str,
	}
	fmt.Fprintf(ctx.Writer, "%s", msg.String())
	flusher.Flush()
}

func EventGetResp(ctx *gin.Context, url string, f func(ctx *gin.Context, id, data []byte) (lastOutPut string, close bool, err error)) (lastOutPut string, err error) {
	client := sse.NewClient(url)
	ch := make(chan *sse.Event)
	err = client.SubscribeChanRawWithContext(ctx, ch)
	if err != nil {
		return "", err
	}
	for {
		select {
		case ed := <-ch:
			if string(ed.Event) == "error" {
				errRender := DefaultRender{}
				_ = json.Unmarshal(ed.Data, &errRender)
				return "", errors.NewError(errRender.ErrNo, errRender.ErrMsg)
			}
			ou, c, errA := f(ctx, ed.ID, ed.Data)
			if errA != nil {
				return "", errA
			}
			if c {
				return ou, nil
			}
		}
	}
	return
}
