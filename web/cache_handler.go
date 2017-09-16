package web

import (
	"github.com/domac/dfc/app"
	"log"
	"net/http"
)

//缓存处理业务
type CacheHandler struct {
}

func NewCacheHandler(f func(ctx *Context)) BaseHandler {
	return BaseHandler{
		Handle: f,
	}
}

//请求缓存
func (self *CacheHandler) Cache(ctx *Context) {
	//获取图片URL
	imageURL := getStringVal("url", ctx.R)

	//不合法的请求
	if !self.checkUrl(imageURL) {
		//拒之门外
		reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "")
		return
	}

	//获取缓存服务
	cacheServer := app.GetCacheServer()
	val, err := cacheServer.Get(imageURL)
	if err != nil {
		log.Printf("[MISS] %s", imageURL)
		maxValLen := 512 * 1024
		val = make([]byte, maxValLen+1)
		err = cacheServer.Set(imageURL, val, 30)
		if err != nil {
			println(err.Error())
		}
	} else {
		log.Printf("[HIT] %s", imageURL)
	}
	reponsePlainText(ctx.W, imageURL)
}

//从伙伴节点的缓存中获取
func (self *CacheHandler) GetFromPeers() (data []byte, err error) {
	return
}

//从本地存储获取
func (self *CacheHandler) GetFromLocal() (data []byte, err error) {
	return
}

//检验URL的合法性
func (self *CacheHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		return false
	}
	return true
}
