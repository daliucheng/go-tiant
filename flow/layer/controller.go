package layer

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/tiant-developer/go-tiant/base"
)

type ControllerErrType string

type IController interface {
	ILayer
	ControllerFunc()
	BindReq() (interface{}, binding.Binding)
	Action() (interface{}, error)
	HandleError(errType ControllerErrType, err error) bool
	ShouldRender() bool

	RenderJsonFail(err error)
	RenderJsonSucc(data interface{})
	RenderJsonAbort(err error)
	RenderJson(code int, message string, data interface{})
}

type Controller struct {
	Layer
}

func (entity *Controller) ControllerFunc() {
	fmt.Print("this is controller func\n")
}

func (entity *Controller) BindReq() (interface{}, binding.Binding) {
	return nil, nil
}

func (entity *Controller) Action() (interface{}, error) {
	return nil, nil
}

func (entity *Controller) HandleError(errType ControllerErrType, err error) bool {
	return false
}

func (entity *Controller) ShouldRender() bool {
	return true
}

func (entity *Controller) RenderJsonFail(err error) {
	base.RenderJsonFail(entity.GetContext(), err)
}

func (entity *Controller) RenderJsonSucc(data interface{}) {
	base.RenderJsonSucc(entity.GetContext(), data)
}

func (entity *Controller) RenderJsonAbort(err error) {
	base.RenderJsonAbort(entity.GetContext(), err)
}

func (entity *Controller) RenderJson(code int, message string, data interface{}) {
	base.RenderJson(entity.GetContext(), code, message, data)
}
