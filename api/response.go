package main

import (
	"encoding/json"
	"io"
	"myproject/video_server/api/defs"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrResponse) {
	// 发送错误码
	w.WriteHeader(errResp.HttpSC)

	// 返回错误详细信息
	resStr, _ := json.Marshal(&errResp.Error)
	// WriteString: 将 string 写入 w 中
	io.WriteString(w, string(resStr))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	// 将 resp 写入到 w 中
	io.WriteString(w, resp)
}
