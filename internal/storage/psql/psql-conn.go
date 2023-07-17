package psql

import (
	"database/sql"
	"log"
	"sync"
)

type psql struct {
	conn *sql.DB
}

var instance *psql = nil
var once sync.Once

func NewPostgresConn(coonectionURI string) *psql {
	once.Do(func() {
		instance = &psql{
			conn: newPostgresConn(coonectionURI),
		}
	})
	return instance
}

func newPostgresConn(coonectionURI string) *sql.DB {
	op := "storage.psql.newPostgresConn"
	conn, err := sql.Open("postgres", coonectionURI)
	if err != nil {
		log.Panicf("%s: %s", op, err)
	}
	return conn
}
