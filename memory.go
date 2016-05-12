// Copyright 2016 orivil Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package session

// memoryStorage for store session in memory. Because "newSession",
// "Destroy" and "Read" function are in the same GC time, and the "Write"
// function didn't used, so it is not necessary to add a locker
type memoryStorage struct {
	sessions map[string]*Session
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		sessions: make(map[string]*Session, 1000),
	}
}

func (sh *memoryStorage) Destroy(ids []string) (destroyed []string) {
	for _, id := range ids {
		delete(sh.sessions, id)
	}
	destroyed = ids
	return
}

func (sh *memoryStorage) Read(id string) *Session {
	if s, ok := sh.sessions[id]; ok {
		return s
	} else {
		s = New(id)
		sh.sessions[id] = s
		return s
	}
}

// GetAll for implement Storage interface
func (sh *memoryStorage) GetAll() (ids []string) { return }

// Write for implement Storage interface
func (sh *memoryStorage) Write(s *Session) {}
