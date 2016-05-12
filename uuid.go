// Copyright 2016 orivil Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package session

import (
	"github.com/satori/go.uuid"
	"net/http"
)

// IdDecoder for generate UUID or read from cookie
type IdDecoder struct {
	cookieKey string
	maxAge    int
}

func NewIdDecoder(cookieKey string, maxAgeMinute int) *IdDecoder {

	return &IdDecoder{
		cookieKey: cookieKey,
		maxAge:    maxAgeMinute * 60,
	}
}

// Read for read id from cookie or generate a new id and it will auto update cookie
func (idc *IdDecoder) Read(w http.ResponseWriter, r *http.Request) (id string) {
	if c, err := r.Cookie(idc.cookieKey); err == http.ErrNoCookie {
		id = uuid.NewV4().String()
	} else if err == nil {
		id = c.Value
	}
	// update cookie max age or send a new cookie
	http.SetCookie(w, &http.Cookie{
		Name:   idc.cookieKey,
		Value:  id,
		Path:   "/",
		MaxAge: idc.maxAge,
	})
	return
}

func (idc *IdDecoder) New(w http.ResponseWriter) (id string) {
	id = uuid.NewV4().String()
	http.SetCookie(w, &http.Cookie{
		Name:   idc.cookieKey,
		Value:  id,
		Path:   "/",
		MaxAge: idc.maxAge,
	})
	return
}