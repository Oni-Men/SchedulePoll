package service

import (
	"log"
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

func TestParseEmbed(t *testing.T) {
	testdata := []struct {
		times  []time.Time
		begins []time.Duration
		ends   []time.Duration
	}{
		{
			times: []time.Time{
				time.Date(2022, 8, 10, 0, 0, 0, 0, time.Local),
				time.Date(2022, 8, 11, 0, 0, 0, 0, time.Local),
				time.Date(2022, 9, 12, 0, 0, 0, 0, time.Local),
			},
			begins: []time.Duration{
				0,
				0,
				0,
			},
			ends: []time.Duration{
				0,
				1 * time.Hour,
				2 * time.Minute,
			},
		},
	}

	for _, td := range testdata {
		p := poll.CreatePoll()
		for i, t := range td.times {
			p.AddColumn(poll.CreateColumn(t, td.begins[i], td.ends[i]))
		}

		embed := PrintPoll(p)
		q, err := ParsePollEmbed(embed)
		if err != nil {
			log.Fatalf("no error expected but we got %v\n", err)
		}

		if p.ID != q.ID {
			t.Fatalf("ID is not the same. %s, %s", p.ID, q.ID)
		}

		if len(p.Columns) != len(q.Columns) {
			t.Fatalf("The length of columns is not the same. %d, %d", len(p.Columns), len(q.Columns))
		}

		if !isColumnsEqual(p, q) {
			t.Fatalf("Columns are not equal.")
		}
	}
}

func isColumnsEqual(a, b *poll.Poll) bool {
	for i, colA := range a.Columns {
		colB := b.Columns[i]
		if colA.Date.Equal(colB.Date) {
			continue
		}

		return false
	}

	return true
}
