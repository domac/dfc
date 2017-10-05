package web

import (
	"errors"
	"fmt"
	"github.com/domac/dfc/app"
	"github.com/domac/dfc/log"
	"github.com/domac/husky"
	"github.com/domac/husky/pb"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var ErrNoPeer = errors.New("no peer found")
var ErrNoPeerData = errors.New("no peer data")

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
			log.GetLogger().Error(err)
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

	ret, err = self.FindCacheDataBySiblings(imageURL)
	if ret != nil && err == nil {
		log.GetLogger().Infof("[SIBLING TCP HIT] %s", imageURL)
		resourceDB.Set([]byte(imageURL), ret)
		return
	}

	//从集群peer中获取
	ret, err = self.FindCacheDataByParents(imageURL)
	if ret != nil && err == nil {
		log.GetLogger().Infof("[PARENT TCP HIT] %s", imageURL)
		resourceDB.Set([]byte(imageURL), ret)
	}
	return
}

//向集群邻居节点请求缓存数据
//TODO
func (self *CacheHandler) FindCacheDataBySiblings(imageURL string) (ret []byte, err error) {
	//获取邻居节点
	siblingNodes, err1 := app.DefaultCacheServer.GetConfig().GetSublingPeerNodes()
	if err1 != nil || siblingNodes == nil {
		return
	}
	wg := sync.WaitGroup{}

	//请求结果
	resultsMap := make(map[string][]byte)
	for _, sp := range siblingNodes {
		wg.Add(1)
		go func(p *app.PeerInfo, lock *sync.WaitGroup) {
			defer lock.Done()
			ret, err := self.getPeerCache(imageURL, p)
			if err == nil && ret != nil {
				log.GetLogger().Infof("[%s] %s get cache success", p.Name, p.Addr)
				resultsMap[p.Addr] = ret
			}
		}(sp, &wg)
	}
	wg.Wait()

	if len(resultsMap) == 0 {
		return nil, ErrNoPeerData
	}

	for addr, ret := range resultsMap {
		if len(ret) > 0 {
			log.GetLogger().Infof("get cache from sibling %s", addr)
			return ret, nil
		}
	}
	return
}

//向集群上层节点请求缓存数据
func (self *CacheHandler) FindCacheDataByParents(imageURL string) (ret []byte, err error) {
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

//获取单个伙伴缓存数据
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
		log.GetLogger().Error(err)
		return nil, err
	}

	if resp != nil {
		bm := &pb.BytesMessage{}
		husky.UnmarshalPbMessage(resp.([]byte), bm)
		return bm.GetBody(), nil
	}

	return nil, nil
}

//缓存key请求检查
func (self *CacheHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		log.GetLogger().Errorf("key [%s] is invalid", url)
		return false
	}
	return true
}

//当前的缓存状态信息
func (self *CacheHandler) CacheStats(ctx *Context) {

	cacheServer, err := app.GetCacheServer()
	if err != nil {
		reponsePlainTextWithStatusCode(ctx.W, http.StatusServiceUnavailable, "")
		return
	}

	entryCount := cacheServer.Cache().EntryCount()
	expiredCount := cacheServer.Cache().ExpiredCount()
	hitRate := cacheServer.Cache().HitRate() //命中率

	resourceDB, err := app.GetResourceDB()
	localKeysCount := len(resourceDB.Keys())
	if err != nil {
		localKeysCount = 0
	}

	//缓存当前信息
	statsInfo := fmt.Sprintf("in_memoty_entry_count : %d\nin_memoty_expired_count: %d\nin_memoty_hitrate: %f\nlocal_keys_count: %d\n",
		entryCount, expiredCount, hitRate, localKeysCount)

	ctx.W.Write([]byte(statsInfo))
}
