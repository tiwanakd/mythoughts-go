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

	stmt := "SELECT content, created, agreecount, disagreecount FROM thoughts ORDER BY created DESC"

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
		err := rows.Scan(&thought.Content, &thought.Created, &thought.AgreeCount, &thought.DisagreeCount)
		if err != nil {
			return nil, err
		}
		thoughts = append(thoughts, thought)
	}

	return thoughts, nil
}
