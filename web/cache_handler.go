package web

import (
	"errors"
	"github.com/domac/dfc/app"
	"github.com/domac/dfc/log"
	"github.com/domac/husky"
	"github.com/domac/husky/pb"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrNoPeer = errors.New("no peer found")

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

	//从in-memory中尝试获取
	val, err := cacheServer.Get(imageURL)
	if err != nil {
		log.GetLogger().Infof("[MEMORY MISS] %s", imageURL)
		b, err := self.FindCacheData(imageURL)
		if err != nil || b == nil || len(b) == 0 {
			reponsePlainTextWithStatusCode(ctx.W, http.StatusBadRequest, "no cache data found")
			return
		}

		err = cacheServer.Set(imageURL, b)
		if err != nil {
			println(err.Error())
		}
		ctx.W.Write(b)
	} else {
		log.GetLogger().Infof("[MEMORY HIT] %s", imageURL)
		ctx.W.Write(val)
	}
}

//用各种方式去获取缓存数据
func (self *CacheHandler) FindCacheData(imageURL string) (ret []byte, err error) {
	//从本地store获取
	resourceDB, err := app.GetResourceDB()
	if err == nil {
		ret, err = resourceDB.Get([]byte(imageURL))
	}

	if ret != nil && err == nil {
		log.GetLogger().Infof("[LOCAL STORE HIT] %s", imageURL)
		return
	}

	log.GetLogger().Infof("[LOCAL STORE MISS] %s", imageURL)

	//从本地目录获取
	ret, err = ioutil.ReadFile(imageURL)
	if ret != nil && err == nil {
		log.GetLogger().Infof("[LOCAL FILE HIT] %s", imageURL)
		resourceDB.Set([]byte(imageURL), ret)
		return
	}
	log.GetLogger().Infof("[LOCAL FILE MISS] %s", imageURL)

	//从集群peer中获取
	ret, err = self.AskPeers(imageURL)
	if ret != nil && err == nil {
		log.GetLogger().Infof("[PEER TCP HIT] %s", imageURL)
		resourceDB.Set([]byte(imageURL), ret)
	}
	return
}

//向集群其它节点请求缓存数据
func (self *CacheHandler) AskPeers(imageURL string) (ret []byte, err error) {
	rr := app.DefaultPeerRoundRobin

	err = ErrNoPeer

	if rr == nil || rr.ParentWrr == nil {
		return
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

	hclient, err := app.CreatePeerSession(p)

	defer func() {
		if hclient != nil {
			log.GetLogger().Infoln("cache handler close hclient")
			hclient.Shutdown()
		}
	}()
	if err != nil {
		log.GetLogger().Infoln("connect to peer fail")
		return nil, err
	}

	req := husky.NewPbBytesPacket(1, "cache_req", []byte(imageURL))
	resp, err := hclient.SyncWrite(*req, 1*time.Second)

	if err != nil {
		println(">>>", err.Error())
		return nil, err
	}

	if resp != nil {
		bm := &pb.BytesMessage{}
		husky.UnmarshalPbMessage(resp.([]byte), bm)
		return bm.GetBody(), nil
	}

	return nil, nil
}

func (self *CacheHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		return false
	}
	return true
}
