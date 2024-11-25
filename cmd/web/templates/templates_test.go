package templates

import (
	"testing"
	"time"

	"github.com/tiwanakd/mythoughts-go/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "PST",
			tm:   time.Date(2024, 11, 11, 17, 15, 0, 0, time.Local),
			want: "11 Nov 2024 at 17:15",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hd := humanDate(test.tm)

			assert.Equal(t, hd, test.want)

		})
	}
}
