package main

import (
	"github.com/julienschmidt/httprouter"
	// "html/temp late"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", homeHandler)
	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)

	router.POST("/api", apiHandler)

	// copy 自 streamserver...
	router.POST("/upload/:vid-id", proxyHandler)

	// statics 指向 template 项目文件夹
	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}
