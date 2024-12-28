package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

func OpenDB(addr string, initfile string) error {

	db, err := sql.Open("sqlite3", addr)

	if err != nil {

		return err
	}

	_db = db

	err = runInitSql(initfile)

	if err != nil {

		_db.Close()

		return fmt.Errorf("failed to run init sql: %v", err)
	}

	return nil
}

func Query(query string, args []any) (*sql.Rows, error) {

	var empty_row *sql.Rows

	results, err := _db.Query(query, args[0:]...)

	if err != nil {

		return empty_row, fmt.Errorf("db query: %s", err.Error())

	}

	return results, err

}

func Exec(query string, args []any) error {

	result, err := _db.Exec(query, args[0:]...)

	if err != nil {
		return fmt.Errorf("db exec: %s", err.Error())
	}

	_, err = result.RowsAffected()

	if err != nil {

		return fmt.Errorf("db exec: rows: %s", err.Error())
	}

	return nil

}

func runInitSql(initfile string) error {

	return nil

}
