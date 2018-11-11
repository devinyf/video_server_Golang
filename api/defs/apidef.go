package defs

// request
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

// response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id`
}

// VideoInfo Video model
type VideoInfo struct {
	ID           string
	AuthorID     int
	Name         string
	DisplayCTime string
}

type Comments struct {
	ID      string
	VideoID string
	Author  string
	Content string
}

type SimpleSession struct {
	UserName string
	TTL      int64
}
