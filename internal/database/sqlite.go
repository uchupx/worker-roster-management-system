package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uchupx/kajian-api/pkg/db"
)

type DBPayload struct {
	Database string
}

// Create new connection to mysql database
// return *sqlx.DB and error
func NewConnection(p DBPayload) (*db.DB, error) {
	conn, err := sqlx.Connect("sqlite3", p.Database)
	if err != nil {
		return nil, fmt.Errorf("failed connectin database host:%s, err: %+v", p.Database, err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed ping database host:%s, err: %+v", p.Database, err)
	}

	return &db.DB{DB: conn}, nil
}
