package db

import "fmt"

type Story struct {
	StoryId          int
	Id               string
	Title            string
	DateMarked       string
	PrimaryMediaName string
	Content          string
}

func SaveStory(id string, title string, dateMarked string, primaryMediaName string, content string) error {

	q := `
	
	INSERT INTO story (
	
		id,
		title,
		date_marked,
		primary_media_name,
		content
	) 
	VALUES (
		?,
		?,
		?,
		?,
		?
	)
	
	`

	a := []any{
		id,
		title,
		dateMarked,
		primaryMediaName,
		content,
	}

	err := exec(q, a)

	if err != nil {
		return fmt.Errorf("failed to save story: %v", err)
	}

	return nil
}

func GetStoryById(id string) (*Story, error) {

	stories := []Story{}

	q := `
	
	SELECT 
		*
	FROM
		story
	WHERE
		id = ?
	`

	a := []any{
		id,
	}

	res, err := query(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get story by id: %v", err)
	}

	defer res.Close()

	for res.Next() {

		story := Story{}

		err = res.Scan(
			&story.StoryId,
			&story.Id,
			&story.Title,
			&story.DateMarked,
			&story.PrimaryMediaName,
			&story.Content,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to get story record: %v", err)
		}

		stories = append(stories, story)

	}

	rlen := len(stories)

	if rlen != 1 {

		return nil, fmt.Errorf("invalid story record len: %d", rlen)
	}

	return &stories[0], nil

}

func GetStoryByTitle(title string) (*Story, error) {

	stories := []Story{}

	q := `
	
	SELECT 
		*
	FROM
		story
	WHERE
		title = ?
	`

	a := []any{
		title,
	}

	res, err := query(q, a)

	if err != nil {

		return nil, fmt.Errorf("failed to get story by title: %v", err)
	}

	defer res.Close()

	for res.Next() {

		story := Story{}

		err = res.Scan(
			&story.StoryId,
			&story.Id,
			&story.Title,
			&story.DateMarked,
			&story.PrimaryMediaName,
			&story.Content,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to get story record: %v", err)
		}

		stories = append(stories, story)

	}

	rlen := len(stories)

	if rlen != 1 {

		return nil, fmt.Errorf("invalid story record len: %d", rlen)
	}

	return &stories[0], nil

}

func DeleteStoryByTitle(title string) (*Story, error) {

	story, err := GetStoryByTitle(title)

	if err != nil {

		return nil, fmt.Errorf("failed to get story by title")
	}

	q := `
	
	DELETE
	FROM
		story
	WHERE
		story_id = ?
	
	`
	a := []any{
		story.StoryId,
	}

	err = exec(q, a)

	if err != nil {
		return nil, fmt.Errorf("failed to delete story record")
	}

	return story, nil

}
