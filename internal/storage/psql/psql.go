package psql

import (
	"database/sql"
	"project/internal/storage"
)

type psql struct {
	conn *sql.DB
}

func (p *psql) CheckDB() error {
	if p.conn == nil {
		return storage.ErrDBNotExists
	}
	return nil
}
