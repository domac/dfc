package web

import (
	"bytes"
	"github.com/domac/dfc/app"
	"io"
	"log"
	"net/http"
	"os"
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

	var buf bytes.Buffer

	if err != nil {
		log.Printf("[MISS] %s", imageURL)
		//测试
		f, err := os.Open(imageURL)
		defer func() {
			if f != nil {
				f.Close()
			}
		}()
		if err != nil {
			log.Println("no file found")
			reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "")
			return
		}

		io.Copy(&buf, f)
		err = cacheServer.Set(imageURL, buf.Bytes(), 30)
		if err != nil {
			println(err.Error())
		}
		ctx.W.Write(buf.Bytes())
	} else {
		log.Printf("[HIT] %s", imageURL)
		ctx.W.Write(val)
	}
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
