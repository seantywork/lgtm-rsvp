package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Comment struct {
	CommentId           int
	Id                  string
	Title               string
	Content             string
	TimestampRegistered string
	TimestampApproved   sql.NullString
}

func ListApprovedComments() ([]Comment, error) {

	var comments = make([]Comment, 0)

	q := `
	
	SELECT
		id,
		title,
		content,
		timestamp_registered
	FROM
		comment
	WHERE
		timestamp_approved IS NOT NULL

	
	`

	a := []any{}

	res, err := query(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get comments: %v", err)
	}

	defer res.Close()

	for res.Next() {

		comment := Comment{}

		err = res.Scan(
			&comment.Id,
			&comment.Title,
			&comment.Content,
			&comment.TimestampRegistered,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to get comment record: %s", err.Error())
		}

		comments = append(comments, comment)

	}

	return comments, nil

}

func GetCommentById(id string) (*Comment, error) {

	var comments = make([]Comment, 0)

	q := `
	
	SELECT
		id,
		title,
		content,
		timestamp_registered
	FROM
		comment
	WHERE
		timestamp_approved IS NULL
		AND id = ?

	
	`

	a := []any{
		id,
	}

	res, err := query(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get comment by id: %v", err)
	}

	defer res.Close()

	for res.Next() {

		comment := Comment{}

		err = res.Scan(
			&comment.Id,
			&comment.Title,
			&comment.Content,
			&comment.TimestampRegistered,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to get comment record by id: %s", err.Error())
		}

		comments = append(comments, comment)

	}

	clen := len(comments)

	if clen != 1 {

		return nil, fmt.Errorf("failed to get comment by id: len: %d", clen)
	}

	return &comments[0], nil

}

func RegisterComment(id string, title string, content string, tr string) error {

	q := `
	
	INSERT INTO comment (
	
		id,
		title,
		content,
		timestamp_registered
	) 
	VALUES (
		?,
		?,
		?,
		?
	)
	
	`

	a := []any{
		id,
		title,
		content,
		tr,
	}

	err := exec(q, a)

	if err != nil {

		return fmt.Errorf("failed to register comment: %v", err)
	}

	return nil
}

func ApproveComment(id string, ta string) error {

	comment, err := GetCommentById(id)

	if err != nil {

		return fmt.Errorf("approve comment: error: %s", id)
	}

	then, err := time.Parse("2006-01-02-15-04-05", comment.TimestampRegistered)

	if err != nil {

		return fmt.Errorf("approve comment: time parse: %s", err.Error())
	}

	now, _ := time.Parse("2006-01-02-15-04-05", ta)

	diff := now.Sub(then)

	if diff.Seconds() > 86400 {

		return fmt.Errorf("approve comment: 24hours passed")
	}

	q := `
	
	UPDATE
		comment
	SET
		timestamp_approved = ?
	WHERE
		id = ?

	`

	a := []any{
		ta,
		id,
	}

	err = exec(q, a)

	if err != nil {

		return fmt.Errorf("faield to approve comment: %v", err)
	}

	return nil
}

func DisapproveCommentByTitle(title string) error {

	comments, err := ListApprovedComments()

	if err != nil {

		return fmt.Errorf("disapprove comment: list approved comments: %s", err.Error())
	}

	for i := 0; i < len(comments); i++ {

		if comments[i].Title == title {
			fmt.Printf("hit: %s\n", comments[i].Id)
			q := `
			
			UPDATE
				comment
			SET
				timestamp_approved = NULL
			WHERE
				id = ?

			`
			a := []any{
				comments[i].Id,
			}

			err = exec(q, a)

			if err != nil {

				return fmt.Errorf("failed to disapprove comment: %v", err)
			}
		}

	}

	return nil
}
