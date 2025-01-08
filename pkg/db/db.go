package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

var initiated bool = false

type SqliteMaster struct {
	TblName string
}

func query(query string, args []any) (*sql.Rows, error) {

	var empty_row *sql.Rows

	results, err := _db.Query(query, args[0:]...)

	if err != nil {

		return empty_row, fmt.Errorf("db query: %s", err.Error())

	}

	return results, err

}

func exec(query string, args []any) error {

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

func OpenDB(addr string) error {

	db, err := sql.Open("sqlite3", addr)

	if err != nil {

		return err
	}

	_db = db

	return nil
}

func Init(initfile string, adminId string, adminPw string) error {

	tables := make([]SqliteMaster, 0)

	admins := make([]Admin, 0)

	q := `
	
	SELECT 
		tbl_name 
	FROM 
		sqlite_master 
	WHERE 
		type='table'
	`

	a := []any{}

	res, err := query(q, a)

	if err != nil {

		return fmt.Errorf("failed to get tables: %v", err)
	}

	defer res.Close()

	for res.Next() {

		t := SqliteMaster{}

		err = res.Scan(&t.TblName)

		if err != nil {

			return fmt.Errorf("failed to read table record: %v", err)

		}

		tables = append(tables, t)

	}

	tlen := len(tables)

	if tlen == 0 {
		log.Printf("no tables found\n")

		err = createFromInitSql(initfile)

		if err != nil {

			return fmt.Errorf("create from init sql: %v", err)
		}

		log.Printf("tables successfully created\n")

	}

	q = `
	
	SELECT
		admin_id
	FROM
		admin
	WHERE
		id = ?

	`
	a = []any{
		adminId,
	}

	res, err = query(q, a)

	if err != nil {

		return fmt.Errorf("failed to get tables: %v", err)
	}

	for res.Next() {

		a := Admin{}

		err = res.Scan(&a.AdminId)

		if err != nil {

			return fmt.Errorf("failed to read admin record: %v", err)

		}

		admins = append(admins, a)

	}

	alen := len(admins)

	if alen == 0 {

		log.Printf("no admins found\n")

		err = addAdmin(adminId, adminPw)

		if err != nil {

			return fmt.Errorf("add admin: %v", err)
		}

		log.Printf("admin successfully added\n")

	}

	initiated = true

	return nil

}

func createFromInitSql(initfile string) error {

	file_b, err := os.ReadFile(initfile)

	if err != nil {

		return fmt.Errorf("failed to read initfile: %v", err)
	}

	q := string(file_b)

	a := []any{}

	err = exec(q, a)

	if err != nil {
		return fmt.Errorf("failed to init: %v", err)
	}

	tables := make([]SqliteMaster, 0)

	q = `
	
	SELECT 
		tbl_name 
	FROM 
		sqlite_master 
	WHERE 
		type='table'
	`

	a = []any{}

	res, err := query(q, a)

	if err != nil {

		return fmt.Errorf("failed to get tables: %v", err)
	}

	defer res.Close()

	for res.Next() {

		t := SqliteMaster{}

		err = res.Scan(&t.TblName)

		if err != nil {

			return fmt.Errorf("failed to read table record: %v", err)

		}

		tables = append(tables, t)

	}

	tlen := len(tables)

	if tlen == 0 {

		return fmt.Errorf("failed to assert tables")
	}

	return nil
}

func addAdmin(id string, pw string) error {

	admins := make([]Admin, 0)

	q := `
	
	INSERT INTO admin (
		id,
		session_id,
		pw
	)
	VALUES (
		?,
		NULL,
		?
	)

	`

	a := []any{
		id,
		pw,
	}

	err := exec(q, a)

	if err != nil {
		return fmt.Errorf("failed to add admin: %v", err)
	}

	q = `
	
	SELECT
		admin_id
	FROM
		admin
	WHERE
		id = ?

	`
	a = []any{
		id,
	}

	res, err := query(q, a)

	if err != nil {

		return fmt.Errorf("failed to add admins: %v", err)
	}

	for res.Next() {

		a := Admin{}

		err = res.Scan(&a.AdminId)

		if err != nil {

			return fmt.Errorf("failed to read admin record: %v", err)

		}

		admins = append(admins, a)

	}

	alen := len(admins)

	if alen == 0 {

		return fmt.Errorf("failed to add admin: 0")
	}

	return nil
}
