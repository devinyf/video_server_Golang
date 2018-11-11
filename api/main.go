package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
}

// RegisterHandlers http.Handler 需要的ServeHTTP方法
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 添加中间件的功能： check session
	ValidateUserSession(r)
	m.r.ServeHTTP(w, r)
}

// 仿构造函数
func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)
	router.POST("/user/:user_name", Login)
	return router
}

func main() {
	r := RegisterHandlers()
	// 中间件劫持
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)

}
