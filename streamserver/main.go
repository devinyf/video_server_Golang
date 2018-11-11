package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// 中间件...劫持http 添加验证
type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 添加的验证机制 判断是否拿到token
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many request!")
		return
	}
	// 被劫持的原始ServeHTTP 方法
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/videos/:vid-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)

	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", mh)
}
