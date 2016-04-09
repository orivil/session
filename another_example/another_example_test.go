package another_example_test

import (
	"fmt"
	. "gopkg.in/orivil/session.v0"
	"net/http"
	"net/url"
	"time"
)

// implement http.CookieJar interface for http client send cookie to server
type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

func ExampleSession() {

	// test set session
	http.HandleFunc("/setsession", func(w http.ResponseWriter, r *http.Request) {

		// 1. new session
		session := NewMemorySession(w, r)

		// 2. use session
		session.Set("login", "foobar")
	})

	// test get session
	http.HandleFunc("/getsession", func(w http.ResponseWriter, r *http.Request) {

		// 1. new session
		session := NewMemorySession(w, r)

		// 2. use session
		value := session.Get("login")

		// print test message, output test result
		fmt.Println(value == "foobar")
	})

	// new http server
	wait := make(chan bool)
	go func() {
		wait <- true
		http.ListenAndServe(":8080", nil)
	}()
	<-wait
	// wait for server to start
	time.Sleep(time.Second)

	// send 3 http requests
	setUrl := "http://localhost:8080/setsession"
	getUrl := "http://localhost:8080/getsession"

	// 1. the first request, this request will build a new session
	response, _ := http.Get(setUrl)
	// store the response cookie(contains session id)
	cookies := response.Cookies()

	// send the cookie(contains session id) to server
	jar := &Jar{cookies: cookies}
	client := http.Client{Jar: jar}
	request, _ := http.NewRequest("GET", getUrl, nil)

	// 2. the second request, this request will get the session which
	// built by the first request, because the client has send the cookie id
	client.Do(request) // output: true (session value == "foobar")

	// 3. the third request, this request will get none session value
	// because the client has no cookie id
	http.Get(getUrl) // output: false (session value == "")

	// Output:
	// true
	// false
}
