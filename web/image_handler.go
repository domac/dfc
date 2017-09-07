package web

import (
	"github.com/gorilla/mux"
	"strings"
)

//图片处理业务
type ImageHandler struct {
}

func NewImageHandler(f func(ctx *Context)) BaseHandler {
	return BaseHandler{
		Handle: f,
	}
}

//图片裁剪
func (self *ImageHandler) Resize(ctx *Context) {
	vars := mux.Vars(ctx.R)
	imageURL := strings.TrimSpace(vars["image"])
	ctx.W.Write([]byte(imageURL))
}
