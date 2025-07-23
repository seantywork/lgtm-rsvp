package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	pkgglob "lgtm-rsvp/pkg/glob"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

var _ismyql bool = false

var initiated bool = false

type SqlMaster struct {
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

func createTLSConf() (*tls.Config, error) {

	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(pkgglob.G_DB_CA_CERT)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca cert: %v", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, fmt.Errorf("failed to add pem")
	}
	clientCert := make([]tls.Certificate, 0, 1)

	certs, err := tls.LoadX509KeyPair(pkgglob.G_DB_CLIENT_CERT, pkgglob.G_DB_CLIENT_KEY)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair")
	}

	clientCert = append(clientCert, certs)

	c := tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	}

	return &c, nil
}

func OpenDB(addr string) error {

	if strings.HasPrefix(addr, pkgglob.G_DB_MYSQL_PREFIX) {

		log.Printf("using mysql\n")

		_ismyql = true

		idPwAddrDb := strings.ReplaceAll(addr, pkgglob.G_DB_MYSQL_PREFIX, "")

		li1 := strings.SplitN(idPwAddrDb, "@", 2)

		if len(li1) != 2 {
			return fmt.Errorf("invalid db info")
		}

		li2 := strings.SplitN(li1[1], "/", 2)

		if len(li2) != 2 {
			return fmt.Errorf("invalid db info")
		}

		connnInfo := fmt.Sprintf("%s@tcp(%s)/%s?tls=custom", li1[0], li2[0], li2[1])

		tc, err := createTLSConf()

		if err != nil {
			return fmt.Errorf("failed to create tls conf: %v", err)
		}
		err = mysql.RegisterTLSConfig("custom", tc)

		if err != nil {
			return fmt.Errorf("failed to register tls conf: %v", err)
		}
		db, err := sql.Open("mysql", connnInfo)

		if err != nil {
			return fmt.Errorf("failed to open mysql connection: %v", err)
		}

		_db = db

	} else {

		log.Printf("using sqlite3\n")

		db, err := sql.Open("sqlite3", addr)

		if err != nil {

			return fmt.Errorf("failed to open sqlite3 connection: %v", err)
		}

		_db = db

	}

	return nil
}

func Init(initfile string, adminId string, adminPw string) error {

	tables := make([]SqlMaster, 0)

	admins := make([]Admin, 0)

	q := ""

	if _ismyql {
		q = `
		
		SELECT 
			table_name 
		FROM 
			information_schema.tables
		`
	} else {
		q = `
		
		SELECT 
			tbl_name 
		FROM 
			sqlite_master 
		WHERE 
			type='table'
		`
	}

	a := []any{}

	res, err := query(q, a)

	if err != nil {

		return fmt.Errorf("failed to get tables: %v", err)
	}

	defer res.Close()

	for res.Next() {

		t := SqlMaster{}

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

	tables := make([]SqlMaster, 0)

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

		t := SqlMaster{}

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
