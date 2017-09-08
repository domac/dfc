package web

import (
	"github.com/domac/dfc/app"
	"net/http"
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
	//获取图片URL
	imageURL := getStringVal("url", ctx.R)

	if !self.checkUrl(imageURL) {
		reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "")
		return
	}

	imageClient := app.GetImgClient()
	err := imageClient.ReadImage(imageURL)
	if err != nil {
		println("resize image error")
		reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "")
		return
	}
	reponsePlainText(ctx.W, imageURL)
}

//检验URL的合法性
func (self *ImageHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		return false
	}
	return true
}
