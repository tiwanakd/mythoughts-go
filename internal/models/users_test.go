package models

import (
	"testing"

	"github.com/tiwanakd/mythoughts-go/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

func TestUserModelAutheticate(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantID   int
		wantErr  error
	}{
		{
			name:     "Valid Credentails",
			email:    "morgan@rdr.com",
			password: "pa$$word",
			wantID:   1,
			wantErr:  nil,
		},
		{
			name:     "Invalid Email",
			email:    "invlid@email.com",
			password: "pa$$word",
			wantID:   0,
			wantErr:  ErrInvalidCredentails,
		},
		{
			name:     "Invalid Password",
			email:    "morgan@rdr.com",
			password: "invalidPassword",
			wantID:   0,
			wantErr:  ErrInvalidCredentails,
		},
		{
			name:     "Invalid Email and Password",
			email:    "invlid@email.com",
			password: "invalidPassword",
			wantID:   0,
			wantErr:  ErrInvalidCredentails,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{DB: db}

			id, err := m.Authenticate(tt.email, tt.password)

			assert.Equal(t, id, tt.wantID)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}
