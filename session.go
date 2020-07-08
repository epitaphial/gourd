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
	for {
		for key, session := range smgr.sessionsMap {
			outdatetime := session.sessionTime.Unix() + int64(smgr.maxLifeTime)
			if time.Now().Unix() > outdatetime {
				delete(smgr.sessionsMap, key)
			}
		}
		time.Sleep(time.Second)
	}
}

func (smgr *SessionManager) setSession(sessionId string, name string, value interface{}) {
	// 服务器中map存储session
	smgr.locker.Lock()
	defer smgr.locker.Unlock()
	if sess, ok := smgr.sessionsMap[sessionId]; ok {
		// 查看原本是否存在该id
		sess.sessionValue[name] = value

	} else {
		session := &Session{
			sessionId:    sessionId,
			sessionTime:  time.Now(),
			sessionValue: make(map[string]interface{}),
		}
		session.sessionValue[name] = value
		smgr.sessionsMap[sessionId] = session
	}
}

func (smgr *SessionManager) getSessionValueBy(sessionId string, sessionName string) (sessionValue interface{}, err error) {
	smgr.locker.RLock()
	defer smgr.locker.RUnlock()
	if sess, ok := smgr.sessionsMap[sessionId]; ok {
		if sessionValue, ok = sess.sessionValue[sessionName]; ok {
			err = nil
		} else {
			err = errors.New("Can not find this session name.")
		}
	} else {
		err = errors.New("Can not find this session.")
	}
	return
}

func (smgr *SessionManager) removeSessionValueBy(sessionId string) (err error) {
	smgr.locker.Lock()
	defer smgr.locker.Unlock()
	if _, ok := smgr.sessionsMap[sessionId]; ok {
		delete(smgr.sessionsMap, sessionId)
	} else {
		err = errors.New("Can not find this session.")
	}
	return
}
