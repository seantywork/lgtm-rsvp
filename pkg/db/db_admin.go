package db

import (
	"database/sql"
	"fmt"
)

type Admin struct {
	AdminId   int
	Id        string
	SessionId sql.NullString
	Pw        sql.NullString
}

func GetAdminById(id string) (*Admin, error) {

	admins := []Admin{}

	q := `
	
	SELECT 
		admin_id,
		id,
		pw
	FROM
		admin
	WHERE
		id = ?
	
	`

	a := []any{
		id,
	}

	res, err := query(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get admin: %v", err)
	}

	defer res.Close()

	for res.Next() {

		admin := Admin{}

		err = res.Scan(
			&admin.AdminId,
			&admin.Id,
			&admin.Pw,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get admin record: %v", err)
		}

		admins = append(admins, admin)

	}

	rlen := len(admins)

	if rlen > 1 {
		return nil, fmt.Errorf("admin record: invalid record len: %d", rlen)
	}

	if rlen == 0 {
		return nil, nil
	}

	return &admins[0], nil

}

func UpsertAdmin(id string, pw string) error {

	admin, err := GetAdminById(id)

	if err != nil {
		return fmt.Errorf("failed to get admin: %v", err)
	}

	q := ""

	a := []any{}

	if admin == nil {

		q = `
		
		INSERT INTO admin (
			id,
			pw
		)
		VALUES (
			?,
			?
		)		
		
		`

		a = append(a, id)
		a = append(a, pw)

	} else {

		q = `
		
		UPDATE
			admin
		SET
			id = ?,
			pw = ?
		WHERE
			admin_id = ?

		`
		a = append(a, id)
		a = append(a, pw)
		a = append(a, admin.AdminId)

	}

	err = exec(q, a)

	if err != nil {
		return fmt.Errorf("failed to upsert admin: %v", err)
	}

	return nil
}

func GetAdminBySessionId(session_id string) (*Admin, error) {

	admins := []Admin{}

	q := `
	
	SELECT 
		admin_id,
		id,
		session_id,
		pw
	FROM
		admin
	WHERE
		session_id = ? 
		AND session_id IS NOT NULL
	
	`

	a := []any{
		session_id,
	}

	res, err := query(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get admin: %v", err)
	}

	defer res.Close()

	for res.Next() {

		admin := Admin{}

		err = res.Scan(
			&admin.AdminId,
			&admin.Id,
			&admin.SessionId,
			&admin.Pw,
		)

		if err != nil {

			return nil, fmt.Errorf("failed to get admin record: %v", err)
		}

		admins = append(admins, admin)

	}

	rlen := len(admins)

	if rlen != 1 {
		return nil, nil
	}

	return &admins[0], nil

}

func SetAdminSessionId(id string, session_id string, signin bool) error {

	admin, err := GetAdminById(id)

	if err != nil {
		return fmt.Errorf("no such admin: %s", id)
	}

	q := ""
	a := []any{}

	if signin {
		q = `
	
		UPDATE
			admin 
		SET
			session_id = ?
		WHERE
			admin_id = ?
	
		`

		a = append(a, session_id)
		a = append(a, admin.AdminId)

	} else {

		q = `
	
		UPDATE
			admin 
		SET
			session_id = NULL
		WHERE
			admin_id = ?
	
		`

		a = append(a, admin.AdminId)

	}

	err = exec(q, a)

	if err != nil {

		return fmt.Errorf("failed to set session: %v", err)
	}

	return nil
}
