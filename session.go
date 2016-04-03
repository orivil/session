package session

import (
	"encoding/json"
)

// Storage to store and recovery the session value
type Storage interface {
	// Destroy for delete the session and return destroied ids
	Destroy(ids []string) (destroied []string)

	// Read for read the data and "Unmarshal" it to Session instance
	// the function "Unmarshal" in this file could be useful or you
	// can use your own function
	Read(id string) *Session

	// Write for "Marshal" session instance to data and write the data
	// to storage
	// the function "Marshal" in this file could be useful or you
	// can use your own function
	Write(s *Session)

	// GetAll for get all of the stored ids, this function only used
	// for init gc time
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

// helper functions
func Unmarshal(data []byte, id string) (s *Session, err error) {
	s = New(id)
	if data != nil {
		err = json.Unmarshal(data, s.Values)
	}
	return s, err
}

func Marshal(s *Session) ([]byte, error) {
	return json.Marshal(s.Values)
}

func CheckChanged(s *Session) bool {
	return s.changed
}

func (s *Session) SetData(key string, data interface{}) {

	s.datas[key] = data
	s.changed = true
}

func (s *Session) GetData(key string) interface{} {

	return s.datas[key]
}

func (s *Session) FlashData(key string) (data interface{}) {
	data = s.datas[key]
	delete(s.datas, key)
	s.changed = true
	return
}

func (s *Session) DelData(key string) {
	delete(s.datas, key)
	s.changed = true
}

func (s *Session) Set(key, value string) {
	s.Values[key] = value
	s.changed = true
}

func (s *Session) Get(key string) (value string) {
	return s.Values[key]
}

// Flash get and delete the value
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
