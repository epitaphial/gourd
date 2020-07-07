package gourd

import (
	"errors"
	"sync"
	"time"
)

type Session struct {
	sessionId    string                 // session对应的id
	sessionTime  time.Time              // session注册的时间
	sessionValue map[string]interface{} // session保存的值,session名：值
}

type SessionManager struct {
	cookieName  string              // 客户端cookie中的名称
	locker      sync.RWMutex        // 读写锁
	maxLifeTime int                 // 单位秒
	sessionsMap map[string]*Session // 保存session的map,键为sessionId
}

func StartSession(cookieName string, maxLifeTime int) *SessionManager {
	smgr := &SessionManager{
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime,
		sessionsMap: make(map[string]*Session),
	}
	go smgr.GC()
	return smgr
}

func (smgr *SessionManager) GC() {

}

func (smgr *SessionManager) setSession(sessionId string, name string, value interface{}) {
	// 服务器中map存储session
	session := &Session{
		sessionId:    sessionId,
		sessionTime:  time.Now(),
		sessionValue: make(map[string]interface{}),
	}
	smgr.locker.Lock()
	defer smgr.locker.Unlock()
	session.sessionValue[name] = value
	smgr.sessionsMap[sessionId] = session
}

func (smgr *SessionManager) getSessionBy(sessionId string) (session *Session, err error) {
	smgr.locker.RLock()
	defer smgr.locker.RUnlock()
	if sess, ok := smgr.sessionsMap[sessionId]; ok {
		err = nil
	} else {
		session = sess
		err = errors.New("can not find this session.")
	}
	return
}
