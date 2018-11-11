package defs

type Err struct {
	ErrorMsg  string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErrResponse{HttpSC: 400, Error: Err{ErrorMsg: "Request body is not correct", ErrorCode: "001"}}
	ErrorNotAuthUser            = ErrResponse{HttpSC: 401, Error: Err{ErrorMsg: "User authentication failed.", ErrorCode: "002"}}
	ErrorDBError                = ErrResponse{HttpSC: 500, Error: Err{ErrorMsg: "DB ops failed", ErrorCode: "003"}}
	ErrorInternalFaults         = ErrResponse{HttpSC: 500, Error: Err{ErrorMsg: "Internal service error", ErrorCode: "004"}}
)
