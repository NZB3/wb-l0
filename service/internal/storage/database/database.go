package database

import (
	"database/sql"
	"log"
)

type db struct {
	connection *sql.DB
}

var instance *db = nil

func NewConnection(connectionURI string) *db {
	if instance == nil {
		instance = &db{
			connection: new(connectionURI),
		}
	}

	return instance
}

func new(connectionURI string) *sql.DB {
	conn, err := sql.Open("postgres", connectionURI)
	if err != nil {
		log.Panicf("Error with connection to database: %s", err)
	}
	return conn
}
