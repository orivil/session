// Copyright 2016 orivil Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package session provide a full-featured session component.
package session

// Storage for storing and recovering the session value
type Storage interface {

	// Destroy for deleting sessions and return destroyed ids
	Destroy(ids []string) (destroyed []string)

	// Read for reading the Session instance
	Read(id string) *Session

	// Write for updating session
	Write(s *Session)

	// GetAll for initializing GC time, it must return all of the stored ids
	GetAll() (ids []string)
}

type Session struct {
	Id      string
	Values  map[string]string
	datas   map[string]interface{}
	changed bool
}

func New(id string) *Session {
	return &Session{
		Id:     id,
		Values: make(map[string]string, 1),
		datas:  make(map[string]interface{}, 1),
	}
}

func (s *Session) ID() string {

	return s.Id
}

func (s *Session) SetData(key string, data interface{}) {

	s.datas[key] = data
}

func (s *Session) GetData(key string) interface{} {

	return s.datas[key]
}

func (s *Session) FlashData(key string) (data interface{}) {
	data = s.datas[key]
	delete(s.datas, key)
	return
}

func (s *Session) DelData(key string) {
	delete(s.datas, key)
}

func (s *Session) Set(key, value string) {
	s.Values[key] = value
	s.changed = true
}

func (s *Session) Get(key string) (value string) {
	return s.Values[key]
}

func (s *Session) Flash(key string) (value string) {
	value = s.Values[key]
	delete(s.Values, key)
	s.changed = true
	return
}

func (s *Session) Del(key string) {
	delete(s.Values, key)
	s.changed = true
}
