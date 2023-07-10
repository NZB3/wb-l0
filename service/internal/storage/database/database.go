package database

import (
	"database/sql"
	"fmt"
)

type db struct {
	connection *sql.DB
}

var instance *db = nil

func NewConnection(connectionURI string) (*db, error) {
	if instance == nil {
		conn, err := new(connectionURI)
		if err != nil {
			return nil, err
		}

		instance = &db{
			connection: conn,
		}
	}

	return instance, nil
}

func new(connectionURI string) (*sql.DB, error) {
	const op = "database.new"

	conn, err := sql.Open("postgres", connectionURI)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return conn, nil
}
