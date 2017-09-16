package web

import (
	"bytes"
	"github.com/domac/dfc/app"
	"io"
	"log"
	"net/http"
	"os"
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

	cacheServer := app.GetCacheServer()
	val, err := cacheServer.Get(imageURL)

	var buf bytes.Buffer

	if err != nil {
		log.Printf("[MISS] %s", imageURL)
		f, err := os.Open(imageURL)
		defer func() {
			if f != nil {
				f.Close()
			}
			buf.Reset()
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

func (self *CacheHandler) GetFromPeers() (data []byte, err error) {
	return
}

func (self *CacheHandler) GetFromLocal() (data []byte, err error) {
	return
}

func (self *CacheHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		return false
	}
	return true
}
