package main

type ApiBody struct {
	Url string `json:"url"`
	Method string `json:"method"`
	ReqBody string `json:"req_body"`
}

type Err struct {
	Error string `json:"error"`
	ErrCode string `json:"err_code:`
}

var (
	ErrorRequestNotRecognized = Err{Error: "api not recognized, bad request", ErrCode: "001"}
	ErrorRequestBodyParseFailed = Err{Error: "request body is not correct", ErrCode: "002"}
	ErrorInternalFaults = Err{Error:"internal service error", ErrCode: "003"}
)