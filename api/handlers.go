package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"myproject/video_server/api/dbops"
	"myproject/video_server/api/defs"
	"myproject/video_server/api/session"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 读取request里的请求体(Body)
	res, _ := ioutil.ReadAll(r.Body)
	fmt.Println("=======res======= :", res)
	ubody := &defs.UserCredential{}
	// 将JSON转化成 go_结构体
	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println(err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionID(ubody.Username)
	su := &defs.SignedUp{Success: true, SessionID: id}

	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}

}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}
