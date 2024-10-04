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
}

type ThoughtModel struct {
	DB *sql.DB
}

func (m *ThoughtModel) ListAll() ([]Thought, error) {
	stmt := "SELECT id, content, created, agreecount, disagreecount FROM thoughts ORDER BY created DESC"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	var thoughts []Thought

	for rows.Next() {
		var thought Thought
		err := rows.Scan(&thought.ID, &thought.Content, &thought.Created, &thought.AgreeCount, &thought.DisagreeCount)
		if err != nil {
			return nil, err
		}
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

func (m *ThoughtModel) Insert(content string) error {
	stmt := "INSERT INTO thoughts (content, created) VALUES ($1, NOW())"
	_, err := m.DB.Exec(stmt, content)
	return err
}
