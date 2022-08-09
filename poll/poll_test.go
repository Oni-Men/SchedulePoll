package poll

import (
	"testing"
)

var count = 0

func createTestVoter() Voter {
	s := string(rune('A' + count))
	count++
	return Voter("Test  Voter " + s)
}

func TestPollVoting(t *testing.T) {
	var err error

	voter1 := createTestVoter()
	voter2 := createTestVoter()
	voter3 := createTestVoter()

	poll := CreatePoll()
	poll.AddColumn("Do you like cats?")
	poll.AddColumn("Do you like dogs?")

	column0, err := poll.GetColumn(0)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}
	if column0.GetContent() != "Do you like cats?" {
		t.Fatalf("column content is wrong, we got %s", column0.GetContent())
	}

	column1, err := poll.GetColumn(1)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}
	if column1.GetContent() != "Do you like dogs?" {
		t.Fatalf("column content is wrong, we got %s", column1.GetContent())
	}

	err = poll.AddVote(voter1, 0)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}

	err = poll.AddVote(voter2, 1)
	if err != nil {
		t.Fatalf("no error expected, but we got %v", err)
	}

	err = poll.AddVote(voter3, 2)
	if err == nil {
		t.Fatalf("error expected, but we got nil")
	}

	err = poll.AddVote(voter3, 1)
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
