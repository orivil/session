// Copyright 2016 orivil Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package session

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"bytes"
	"strings"
)

type DBSession struct {
	ID     int `gorm:"primary_key"`
	SID    string `gorm:"type:char(36);index;not null;"`
	Values string `gorm:"type:varchar(255);not null"`
}


type DBStorage struct {
	db *gorm.DB
	// session cache
	sessions map[string]*Session
}

func NewDBStorage(driver string, db *sql.DB) *DBStorage {

	gormDB, err := gorm.Open(driver, db)
	if err != nil {
		panic(err)
	}
	gormDB.AutoMigrate(&DBSession{})
	return &DBStorage{db:gormDB}
}

func (this *DBStorage) Destroy(ids []string) (destroyed []string) {

	delNum := this.db.Delete(DBSession{}, "s_id in (?)", ids).RowsAffected
	if delNum == int64(len(ids)) {

		return ids
	} else {
		var unDestroyed []string
		this.db.Model(&DBSession{}).Where("s_id in (?)", ids).Pluck("s_id", &unDestroyed)
		for _, id := range ids {
			delete := true
			for _, exist := range unDestroyed {
				if id == exist {
					delete = false
					break
				}
			}
			if delete {
				destroyed = append(destroyed, id)
			}
		}
	}

	for _, id := range ids {
		delete(this.sessions, id)
	}
	return
}

func (this *DBStorage) Read(id string) *Session {

	if session, ok := this.sessions[id]; ok {

		return session
	} else {
		dbSession := &DBSession{SID: id}
		this.db.FirstOrCreate(dbSession, "s_id=?", id)
		session := New(id)
		if len(dbSession.Values) > 0 {
			values := bytes.Split([]byte(dbSession.Values), []byte("|"))
			vlength := len(values)
			for i := 0; i < vlength; i += 2 {
				session.Values[string(values[i])] = string(values[i + 1])
			}
		}

		this.sessions[id] = session
		return session
	}
}

func (this *DBStorage) Write(s *Session) {
	if len(s.Values) > 0 {
		values := bytes.NewBuffer(nil)
		for key, value := range s.Values {
			values.WriteString(key)
			values.WriteString("|")
			values.WriteString(value)
			values.WriteString("|")
		}

		this.db.Model(&DBSession{}).Where("s_id = ?", s.Id).
		Update("values", string(values.Bytes()[0:values.Len() - 2]))
	}
}

func (this *DBStorage) GetAll() (ids []string) {
	var sessions []*DBSession
	this.db.Find(&sessions)
	ids = make([]string, len(sessions))
	this.sessions = make(map[string]*Session, len(sessions))
	for idx, dbs := range sessions {
		ids[idx] = dbs.SID
		this.sessions[dbs.SID] = trans(dbs)
	}
	return
}

func trans(dbs *DBSession) (s *Session) {
	s = &Session{
		Id: dbs.SID,
		Values: make(map[string]string, 1),
	}
	if len(dbs.Values) > 0 {
		values := strings.Split(dbs.Values, "|")
		vlength := len(values)
		for i := 0; i < vlength; i += 2 {
			s.Values[values[i]] = values[i + 1]
		}
	}
	return
}