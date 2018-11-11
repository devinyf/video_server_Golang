/*
httpclient 代理发送真正的 api_request
*/
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	switch b.Method {
	case http.MethodGet:
		// get 方法没有body 填 nil
		// http.NewRequest就是直接调用后端接口，
		//其中第二个参数会从配置里读取LB的address和port，自然就会转发到8000那个api端口。
		req, _ := http.NewRequest("GET", b.Url, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)

	case http.MethodPost:
		// POST 方法
		req, _ := http.NewRequest("POST", b.Url, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)

	case http.MethodDelete:
		// Delete 方法
		req, _ := http.NewRequest("GET", b.Url, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request !")
		return
	}
}

// 将后端返回的Response再返回
func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// 说明服务器没有响应
		re, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(re))
		return
	}

	// 转发状态吗
	w.WriteHeader(r.StatusCode)
	// 转发应答体
	io.WriteString(w, string(res))
}
