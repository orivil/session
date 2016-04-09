package session_test

import (
	. "gopkg.in/orivil/session.v0"
	"net/http"
)

func ExampleMemorySession() {

	//use memory to sotre session
	http.HandleFunc("/memory", func(w http.ResponseWriter, r *http.Request) {

		// 1. new session
		session := NewMemorySession(w, r)

		// 2. use session
		session.Set("login", "true")
	})
}

func ExamplePermanentSession() {

	// use long time storage to store session
	// You must implement the "session.Storage" interface by yourself.
	var storage Storage

	// 1. set sorage
	SetStorage(storage)

	http.HandleFunc("/storage", func(w http.ResponseWriter, r *http.Request) {

		// 2. new session
		session := NewPermanentSession(w, r)

		// 3. use session
		session.Set("login", "true")

		// 4. store session
		StorePermanentSession(session)
	})
}
