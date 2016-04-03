package session

// memoryStorage for store session in memory, Because
type memoryStorage struct {
	sessions map[string]*Session
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		sessions: make(map[string]*Session, 1000),
	}
}

// Because "newSession", "Destroy" and "Read" function are in the same GC time, and the "Write"
// function didn't used, so this is not necessary to add a locker
func (sh *memoryStorage) Destroy(ids []string) (destroied []string) {
	for _, id := range ids {
		delete(sh.sessions, id)
	}
	destroied = ids
	return
}

func (sh *memoryStorage) newSession(id string) *Session {
	s := New(id)
	sh.sessions[id] = s
	return s
}

func (sh *memoryStorage) Read(id string) *Session {

	return sh.sessions[id]
}

// GetAll for implement Storage interface, it's useless for memory storage
func (sh *memoryStorage) GetAll() (ids []string) { return }

// Write for implement Storage interface, it's useless for memory storage
// The session was stored when it be created, see the "Read" function
func (sh *memoryStorage) Write(s *Session) {}
