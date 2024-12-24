package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

func OpenDB(addr string) error {

	db, err := sql.Open("sqlite3", addr)

	if err != nil {

		return err
	}

	_db = db

	return nil
}
