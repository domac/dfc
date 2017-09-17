package web

import (
	"github.com/domac/dfc/app"
	"io/ioutil"
	"log"
	"net/http"
)

//cache api handler
type CacheHandler struct {
}

func NewCacheHandler(f func(ctx *Context)) BaseHandler {
	return BaseHandler{
		Handle: f,
	}
}

func (self *CacheHandler) Cache(ctx *Context) {
	imageURL := getStringVal("url", ctx.R)

	if !self.checkUrl(imageURL) {
		//reject
		reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "")
		return
	}

	cacheServer, err := app.GetCacheServer()

	if err != nil {
		reponsePlainTextWithStatusCode(ctx.W, http.StatusServiceUnavailable, "")
		return
	}

	val, err := cacheServer.Get(imageURL)

	if err != nil {
		log.Printf("[MISS] %s", imageURL)
		//读取文件数据
		b, err := ioutil.ReadFile(imageURL)
		if err != nil {
			log.Println("no file found")
			reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "no file found")
			return
		}
		err = cacheServer.Set(imageURL, b)
		if err != nil {
			println(err.Error())
		}
		ctx.W.Write(b)
	} else {
		log.Printf("[HIT] %s", imageURL)
		ctx.W.Write(val)
	}
}

func (self *CacheHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		return false
	}
	return true
}
