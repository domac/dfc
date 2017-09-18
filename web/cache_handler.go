package web

import (
	"github.com/domac/dfc/app"
	"github.com/domac/husky"
	"github.com/domac/husky/pb"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//cache api handler
type CacheHandler struct {
}

func NewCacheHandler(f func(ctx *Context)) BaseHandler {
	return BaseHandler{
		Handle: f,
	}
}

//请求缓存
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
			//向集群节点询问
			ret, err := self.AskPeers(imageURL)
			if err != nil || ret == nil || len(ret) == 0 {
				reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "no resource found")
				return
			}
			b = ret
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

//向集群其它节点请求缓存数据
func (self *CacheHandler) AskPeers(imageURL string) (ret []byte, err error) {
	rr := app.DefaultPeerRoundRobin

	//没有集群信息
	if len(rr.PeersSessions) == 0 {
		return nil, nil
	}

	if rr.ParentWrr == nil {
		return nil, nil
	}
	//获取合适权重的节点
	np := rr.ParentWrr.Next()
	if np != nil {
		p := np.(*app.PeerInfo)
		ret, err = self.getPeerCache(imageURL, p)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (self *CacheHandler) getPeerCache(imageURL string, p *app.PeerInfo) ([]byte, error) {

	if hc, ok := app.DefaultPeerRoundRobin.PeersSessions["parent"]; ok {
		if hc != nil {
			//响应请求
			req := husky.NewPbBytesPacket(1, "cache_req", []byte(imageURL))
			resp, err := hc.SyncWrite(*req, 500*time.Millisecond)
			if err != nil {
				return nil, err
			}
			if resp != nil {
				bm := &pb.BytesMessage{}
				husky.UnmarshalPbMessage(resp.([]byte), bm)
				return bm.GetBody(), nil
			}
		}
	}
	return nil, nil
}

func (self *CacheHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		return false
	}
	return true
}
