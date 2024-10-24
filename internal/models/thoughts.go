package models

import (
	"database/sql"
	"errors"
	"time"
)

type Thought struct {
	ID            int
	Content       string
	Created       time.Time
	AgreeCount    int
	DisagreeCount int
	UserID        int
	Username      string
}

type ThoughtModel struct {
	DB *sql.DB
}

func (m *ThoughtModel) List(sortby string) ([]Thought, error) {
	var stmt string
	if sortby == "created" {
		stmt = "SELECT id, content, created, agreecount, disagreecount, user_id FROM thoughts ORDER BY created DESC"
	} else if sortby == "agree" {
		stmt = "SELECT id, content, created, agreecount, disagreecount, user_id FROM thoughts ORDER BY agreecount DESC"
	} else if sortby == "disagree" {
		stmt = "SELECT id, content, created, agreecount, disagreecount, user_id FROM thoughts ORDER BY disagreecount DESC"
	}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	var thoughts []Thought

	//create a new variable that holds the UserModel with DB provided by ThoughtModel
	//this needs to be used to fetch the username which will be added to Thoughts slice
	users := UserModel{m.DB}

	for rows.Next() {
		var thought Thought
		err := rows.Scan(&thought.ID, &thought.Content, &thought.Created, &thought.AgreeCount, &thought.DisagreeCount, &thought.UserID)
		if err != nil {
			return nil, err
		}

		user, err := users.Get(thought.UserID)
		if err != nil {
			return nil, err
		}

		thought.Username = user.Username
		thoughts = append(thoughts, thought)
	}

	return thoughts, nil
}

func (m *ThoughtModel) AddLike(id int) (int, error) {
	query := "UPDATE thoughts SET agreecount = agreecount + 1 WHERE id = $1 RETURNING agreecount"

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var agreeCount int
	err = stmt.QueryRow(id).Scan(&agreeCount)
	if err != nil {
		return 0, err
	}

	return agreeCount, nil
}

func (m *ThoughtModel) AddDislike(id int) (int, error) {
	query := "UPDATE thoughts SET disagreecount = disagreecount + 1 WHERE id = $1 RETURNING disagreecount"

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var disagreeCount int
	err = stmt.QueryRow(id).Scan(&disagreeCount)
	if err != nil {
		return 0, err
	}

	return disagreeCount, nil
}

func (m *ThoughtModel) Insert(content string, userID int) (Thought, error) {
	query := `INSERT INTO thoughts (content, created, user_id) VALUES ($1, NOW(), $2)
	RETURNING id, content, created, user_id`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return Thought{}, err
	}

	var thought Thought
	err = stmt.QueryRow(content, userID).Scan(&thought.ID, &thought.Content, &thought.Created, &thought.UserID)
	if err != nil {
		return Thought{}, err
	}

	return thought, nil
}
