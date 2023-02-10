package poll

import (
	"testing"

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
