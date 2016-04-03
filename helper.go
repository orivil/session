package session

import (
	"log"
	"net/http"
	"time"
)

var MemoryGC *SessionGC
var PermanentGC *SessionGC

func init() {
	mStorage := newMemoryStorage()
	memoryConstructor := func(id string) *Session {

		return mStorage.newSession(id)
	}

	permanentConstructor := func(id string) *Session {

		return New(id)
	}

	MemoryGC = NewSessionGC(memoryConstructor)
	MemoryGC.SetStorage(mStorage)

	PermanentGC = NewSessionGC(permanentConstructor)

	// set default options
	maxAge := 45
	checkNum := 3
	MemoryGC.Config(maxAge, "orivil-memory-session", checkNum)
	PermanentGC.Config(maxAge, "orivil-storage-session", checkNum)
}

func ConfigMemory(maxAgeMinute, checkNum int, cookieKey string) {
	MemoryGC.Config(maxAgeMinute, cookieKey, checkNum)
}

func ConfigPermanent(maxAgeMinute, checkNum int, cookieKey string) {
	PermanentGC.Config(maxAgeMinute, cookieKey, checkNum)
}

func SetStorage(s Storage) {

	PermanentGC.SetStorage(s)
}

func NewMemorySession(w http.ResponseWriter, r *http.Request) *Session {

	return MemoryGC.Read(w, r)
}

func NewPermanentSession(w http.ResponseWriter, r *http.Request) *Session {

	return PermanentGC.Read(w, r)
}

func StorePermanentSession(s *Session) {
	if CheckChanged(s) {
		PermanentGC.storage.Write(s)
	}
}

func Log() {
	printmsg(MemoryGC, "Memory GC")
	printmsg(PermanentGC, "Permanent GC")
}

func printmsg(gc *SessionGC, name string) {
	expired := 0
	now := time.Now()
	gc.mu.Lock()
	defer gc.mu.Unlock()

	for _, t := range gc.times {
		if now.After(t) {
			expired++
		}
	}
	log.Printf("%s session left: %d, expired: %d\n", name, len(gc.times), expired)
}
