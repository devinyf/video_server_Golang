package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

// 将服务器的视频文件 stream 到客户端
func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid

	video, err := os.Open(vl)
	if err != nil {
		log.Printf("Error when try to open file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error!")
		return
	}
	defer video.Close()

	// 指定 mp4 格式， 浏览器会以mp4解析播放
	w.Header().Set("Content-Type", "video/mp4")
	// 将内容以二进制流传输给 client 端
	http.ServeContent(w, r, "", time.Now(), video)
}

// 将 client 端的视频上传到 server
func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// MaxBytesReader 限定可读最大文件大小
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big !!")
		return
	}

	// FormFile: 从前端 HTML 的 form 表单中取出的name <form name="file">
	file, _, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error !!!")
		return
	}
	defer file.Close()

	// 二进制读取 file 返回 byte Array(这里是一个视频文件)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error : %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error !!")
	}
	// 从url中提取 视频 的名子
	fn := p.ByName("vid-id")
	// 写入文件到服务器 VIDEO_DIR 文件夹
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("Write file err : %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error !!!")
		return
	}
	// 返回状态吗
	w.WriteHeader(http.StatusCreated)
	// 返回给浏览器提示：
	io.WriteString(w, "Upload successfully !")
}
