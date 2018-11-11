package session

import (
	"myproject/video_server/api/dbops"
	"myproject/video_server/api/defs"
	"sync"
	"time"

	"github.com/pborman/uuid"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

// nowInMilli - 当前时间 毫秒
func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

func LoadSessionFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss) //sync.Map 的函数
		return true

	})
}

func GenerateNewSessionID(un string) string {
	id := uuid.New()
	ct := nowInMilli()     // 当前时间-毫秒
	ttl := ct + 30*60*1000 // 过期时间:30分钟

	ss := &defs.SimpleSession{UserName: un, TTL: ttl}
	// 储存在sessionMap中
	sessionMap.Store(id, ss)
	// 储存在database中
	dbops.InsertSession(id, ttl, un)

	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*defs.SimpleSession).UserName, false
	}
	return "", true
}
