package service

import (
	"testing"
	"time"

	"github.com/Oni-Men/SchedulePoll/internal/model/poll"
)

func TestReactionValidator(t *testing.T) {
	p := poll.CreatePoll()
	p.AddColumn(poll.CreateColumn(time.Now(), 0, 0))
	p.AddColumn(poll.CreateColumn(time.Now(), 0, 0))
	p.AddColumn(poll.CreateColumn(time.Now(), 0, 0))

	testdata := []struct {
		id       string
		expected bool
	}{
		{
			"\U0001F1E5",
			false,
		}, {
			"🇦",
			true,
		}, {
			"🇨",
			true,
		}, {
			"🇩",
			false,
		},
	}

	for _, td := range testdata {
		if isValidReaction(td.id, p) != td.expected {
			t.Fatalf("failed to validate for %s.", td.id)
		}
	}
}
