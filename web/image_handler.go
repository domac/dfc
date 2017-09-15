package web

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
	reponsePlainText(ctx.W, imageURL)
}

//检验URL的合法性
func (self *CacheHandler) checkUrl(url string) bool {
	if url == "" || len(url) < 5 {
		return false
	}
	return true
}
