// Copyright 2016 orivil Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"sync"
	"time"
)


type SessionGC struct {

	storage         Storage
	session_max_age time.Duration
	uuid            *IdDecoder
	checkNum        int

	// time is the expired time
	times map[string]time.Time
	mu    sync.RWMutex
}

func NewSessionGC() *SessionGC {
	return &SessionGC{}
}

func (gc *SessionGC) SetStorage(s Storage) {

	gc.storage = s

	// set all of the stored session's gc time to "now"
	ids := gc.storage.GetAll()
	gc.times = make(map[string]time.Time, len(ids))
	now := time.Now()
	for _, id := range ids {
		gc.times[id] = now.Add(gc.session_max_age)
	}
}

func (gc *SessionGC) Config(maxAgeMinute int, cookieKey string, checkNum int) {
	gc.session_max_age = time.Duration(maxAgeMinute) * time.Minute
	gc.uuid = NewIdDecoder(cookieKey, maxAgeMinute)
	gc.checkNum = checkNum
}

func (gc *SessionGC) GetID(w http.ResponseWriter, r *http.Request) string {
	sessionID := gc.uuid.Read(w, r)
	now := time.Now()
	gc.mu.RLock()
	t, ok := gc.times[sessionID]
	gc.mu.RUnlock()

	// session 已过期
	if ok && now.After(t) {
		// TODO： 是否应当将当前过期的数据删除, 大多数情况是可以继续使用的, 还可以提升用户体验
		// 更新 session id
		sessionID = gc.uuid.New(w)
	}
	return sessionID
}

func (gc *SessionGC) Destroy(now time.Time) {
	// 每次随机检查 n 条记录是否过期, 当过期的记录比例越多的时候, 则被删除的概率就越大
	checkNum := gc.checkNum
	// 记录过期的数据
	var expireds []string
	for id, t := range gc.times {
		if now.After(t) {
			expireds = append(expireds, id)
		}
		if checkNum <= 1 {
			break
		}
		checkNum--
	}

	// 删除 gc 中的记录
	if len(expireds) > 0 {
		// 在 gc 中删除被 storage 成功删除的记录
		destroyed := gc.storage.Destroy(expireds)
		for _, delID := range destroyed {
			delete(gc.times, delID)
		}
	}
}

func (gc *SessionGC) Read(sessionID string) (s *Session) {
	
	now := time.Now()
	gc.mu.Lock()
	defer gc.mu.Unlock()

	s = gc.storage.Read(sessionID)

	// 更新当前记录
	gc.times[sessionID] = now.Add(gc.session_max_age)

	gc.Destroy(now)
	return
}
