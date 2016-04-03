package session

import (
	"net/http"
	"sync"
	"time"
)

type SessionConstructor func(sessionID string) *Session

type SessionGC struct {

	// Because the memory only stored the "Session" ptr, so it is not necessary
	// to cover the session data after used, but memory storage have to store
	// the "Session" ptr when it created, and this step was did by GC, so the
	// GC have to open up the "Session" constructor for memory storage
	constructor SessionConstructor

	storage         Storage
	session_max_age time.Duration
	uuid            *IdDecoder
	checkNum        int

	// time is the expired time
	times map[string]time.Time
	mu    sync.Mutex
}

func NewSessionGC(c SessionConstructor) *SessionGC {
	return &SessionGC{constructor: c}
}

func (gc *SessionGC) SetStorage(s Storage) {

	gc.storage = s

	// set all of the stored session's gc time to "now"
	ids := gc.storage.GetAll()
	gc.times = make(map[string]time.Time, len(ids))
	now := time.Now()
	for _, id := range ids {
		gc.times[id] = now
	}
}

func (gc *SessionGC) Config(maxAgeMinute int, cookieKey string, checkNum int) {
	gc.session_max_age = time.Duration(maxAgeMinute) * time.Minute
	gc.uuid = NewIdDecoder(cookieKey, maxAgeMinute)
	gc.checkNum = checkNum
}

func (gc *SessionGC) Read(w http.ResponseWriter, r *http.Request) (s *Session) {
	// 每次随机检查 n 条记录是否过期, 当过期的记录比例越多的时候, 则被删除的概率就越大
	checkNum := gc.checkNum

	// 读取或生成 id
	sessionID := gc.uuid.Read(w, r)

	// 记录过期的数据
	var expireds []string
	now := time.Now()
	gc.mu.Lock()
	defer gc.mu.Unlock()
	// 随机查询 session 是否过期
	for id, t := range gc.times {
		if now.After(t) {
			expireds = append(expireds, id)
		}
		if checkNum <= 1 {
			break
		}
		checkNum--
	}

	// 查询当前 session id 是否有记录及是否过期
	if t, ok := gc.times[sessionID]; ok && now.After(t) {
		// 有记录但已过期

		// TODO： 是否应当将当前过期的数据删除, 为何不可以继续使用？
		// session gc 的作用应当是避免 session 数据无限膨胀
		expireds = append(expireds, sessionID)
	} else if ok {
		// 有记录且没过期
		s = gc.storage.Read(sessionID)
	}

	// 删除 gc 中的记录
	if len(expireds) > 0 {
		// 在 gc 中删除被 storage 成功删除的记录
		destroied := gc.storage.Destroy(expireds)
		for _, delID := range destroied {
			delete(gc.times, delID)
		}
	}

	// 如果没有读到 session 数据, 则构造一个新 session
	if s == nil {
		s = gc.constructor(sessionID)
	}

	// 更新当前记录
	gc.times[sessionID] = now.Add(gc.session_max_age)
	return
}
