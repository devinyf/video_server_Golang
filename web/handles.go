package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

// 替换 home.html 的信息
type HomePage struct {
	Name string
}

// 替换 user.html 的信息
type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 通过cookie判断是否登陆
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		p := &HomePage{Name: "squall"}
		// 指定需要解析的 html 文件路径
		t, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			log.Printf("Parsing template home.html error: %s", err)
			return
		}
		// 将p 通过 w(ResponseWrite) 返回给前端
		t.Execute(w, p)
		return
	}

	// 如果登陆过， 跳转到 userhome 页面
	// 这里只判断是否存在, 前端 js 通过sessionID 判断是用户名否匹配
	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")
	// 如果没有登陆, 重定向 302
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 从 request 登陆页提交的表单中提取 username, 如值不存在返回 ""
	// 用户是否合法 在前端调用后端 API 来验证
	fname := r.FormValue("username")

	var p *UserPage
	// cookie 中存在数据
	if len(cname.Value) != 0 {
		// 从 cookie 中取出 username
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		// 从表单提交的信息中读取
		p = &UserPage{Name: fname}
	}
	t, e := template.ParseFiles("./templates/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error %s", e)
		return
	}
	// 提交渲染
	t.Execute(w, p)

}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(re))
		return
	}

	request(apibody, w, r)
	defer r.Body.Close()
}

// proxy 代理 转发
func proxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse("http://127.0.0.1:9000/")
	// 通过代理 将域名替换为 u, URL_path不变
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
