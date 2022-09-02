package poll

import (
	"log"
	"testing"
	"time"

	"github.com/Oni-Men/SchedulePoll/pkg/dateparser"
)

func TestPollVoting(t *testing.T) {
	var err error

	poll := CreatePoll()
	parser := dateparser.NewDateParser("2022/4/1,2022/4/2")

	for parser.HasNext() {
		res, err := parser.Next()
		if err != nil {
			t.Fatalf("no error expected. but we got %v\n", err)
		}
		poll.AddColumn(CreateColumn(res.Date, res.BeginAt, res.EndAt))
	}

	// poll.AddColumn("Do you like cats?")
	// poll.AddColumn("Do you like dogs?")

	err = poll.AddVote(0)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}

	err = poll.AddVote(1)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}

	err = poll.AddVote(2)
	if err == nil {
		t.Fatalf("error expected, but we got nil")
	}

	err = poll.AddVote(1)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}

	votes, err := poll.GetVotes(0)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}
	if votes != 1 {
		t.Fatalf("1 expected, but we got %d", votes)
	}

	votes, err = poll.GetVotes(1)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}
	if votes != 2 {
		t.Fatalf("2 expected, but we got %d", votes)
	}

	_, err = poll.GetVotes(2)
	if err == nil {
		t.Fatalf("error expected, but we got nil")
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
		p := CreatePoll()
		for i, t := range td.times {
			p.AddColumn(CreateColumn(t, td.begins[i], td.ends[i]))
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

func isColumnsEqual(a, b *Poll) bool {
	for i, colA := range a.Columns {
		colB := b.Columns[i]
		if colA.Date.Equal(colB.Date) {
			continue
		}

		return false
	}

	return true
}
