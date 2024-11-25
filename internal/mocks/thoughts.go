package mocks

import (
	"time"

	"github.com/tiwanakd/mythoughts-go/internal/models"
)

var mockThought1 = models.Thought{
	ID:            1,
	Content:       "You will come over trust me... All this will pass.",
	Created:       time.Now(),
	AgreeCount:    3,
	DisagreeCount: 0,
	UserID:        1,
	Username:      "test",
}

var mockThought2 = models.Thought{
	ID:            2,
	Content:       "You need to forgvie yourself dude... There is a lot more to life.",
	Created:       time.Now(),
	AgreeCount:    1,
	DisagreeCount: 5,
	UserID:        1,
	Username:      "test",
}

type ThoughtModel struct{}

func (m *ThoughtModel) GetThoughts(stmt string, args ...any) ([]models.Thought, error) {
	return []models.Thought{mockThought1, mockThought2}, nil
}

func (m *ThoughtModel) List(sortby string) ([]models.Thought, error) {
	switch sortby {
	case "agree":
		return []models.Thought{mockThought1, mockThought2}, nil
	case "disagree":
		return []models.Thought{mockThought2, mockThought1}, nil
	default:
		return []models.Thought{mockThought1, mockThought2}, nil
	}
}

func (m *ThoughtModel) AddLike(id int) (int, error) {
	switch id {
	case 1:
		return 4, nil
	case 2:
		return 2, nil
	default:
		return 0, models.ErrNoRecord
	}
}

func (m *ThoughtModel) AddDislike(id int) (int, error) {
	switch id {
	case 1:
		return 1, nil
	case 2:
		return 6, nil
	default:
		return 0, models.ErrNoRecord
	}
}

func (m *ThoughtModel) Insert(content string, userID int) (models.Thought, error) {
	newThought := models.Thought{
		ID:            3,
		Content:       "This is a test post from a test Application and Test Server.",
		Created:       time.Now(),
		AgreeCount:    0,
		DisagreeCount: 0,
		UserID:        1,
		Username:      "test",
	}

	if userID == 1 {
		return newThought, nil
	}

	return models.Thought{}, models.ErrInvalidCredentails
}

func (m *ThoughtModel) UserThoughts(userID int, sortby string) ([]models.Thought, error) {
	switch userID {
	case 1:
		return []models.Thought{mockThought1}, nil
	case 2:
		return []models.Thought{mockThought2}, nil
	default:
		return []models.Thought{}, nil
	}
}

func (m *ThoughtModel) DeleteThought(id int) error {
	if id == 1 {
		return nil
	}
	return models.ErrNoRecord
}
