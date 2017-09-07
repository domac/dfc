package web

import (
	"net/http"
)

func InitServer() (*http.Server, error) {
	return loadRouters()
}

//上下文信息
type Context struct {
	R         *http.Request
	W         http.ResponseWriter
	callbacks []func()
	Data      map[interface{}]interface{}
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		R:         r,
		W:         w,
		callbacks: make([]func(), 0, 2),
		Data:      make(map[interface{}]interface{}, 2),
	}
}

//方法回调
func (ctx *Context) CallBack(f func()) {
	ctx.callbacks = append(ctx.callbacks, f)
}

func (ctx *Context) Done() {
	n := len(ctx.callbacks) - 1
	for i := n; i >= 0; i-- {
		ctx.callbacks[i]()
	}
}

//基础服务处理类
type BaseHandler struct {
	Ctx    map[string]interface{}
	Handle func(ctx *Context)
}

func NewBaseHandler(f func(ctx *Context)) BaseHandler {
	return BaseHandler{
		Handle: f,
	}
}

//http服务处理
func (b BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err_ := recover()
		if err_ == nil {
			return
		}
		return
	}()

	ctx := newContext(w, r)
	defer ctx.Done()
	b.Handle(ctx)
}
