// Copyright 2016 orivil Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package session

import (
	"log"
	"net/http"
)

var MemoryGC *SessionGC
var PermanentGC *SessionGC

func init() {
	mStorage := newMemoryStorage()

	MemoryGC = NewSessionGC()
	MemoryGC.SetStorage(mStorage)

	PermanentGC = NewSessionGC()

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
	sessionID := MemoryGC.GetID(w, r)
	return MemoryGC.Read(sessionID)
}

func NewPermanentSession(w http.ResponseWriter, r *http.Request) *Session {

	sessionID := PermanentGC.GetID(w, r)
	return PermanentGC.Read(sessionID)
}

func StorePermanentSession(s *Session) {
	if s.changed {
		s.changed = false
		PermanentGC.storage.Write(s)
	}
}

func Log() {
	printmsg(MemoryGC, "Memory GC")
	printmsg(PermanentGC, "Permanent GC")
}

func printmsg(gc *SessionGC, name string) {

	log.Printf("%s session left: %d\n", name, len(gc.times))
}
