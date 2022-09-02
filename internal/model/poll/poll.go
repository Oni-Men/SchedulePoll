package poll

import (
	"errors"
	"log"
	"time"

	"github.com/Oni-Men/SchedulePoll/pkg/rands"
)

type Poll struct {
	ID          string
	Title       string
	Description string
	Columns     []*Column
	Due         time.Time
}

func CreatePoll() *Poll {
	id, err := rands.RandomHex(3)
	if err != nil {
		log.Fatalf("failed to create a new poll: %v\n", err)
	}
	p := Poll{
		ID:      id,
		Columns: make([]*Column, 0, 10),
	}
	return &p
}

func (p *Poll) AddColumnsAll(list []*Column) error {
	for _, t := range list {
		if _, err := p.AddColumn(t); err != nil {
			return err
		}
	}
	return nil
}

func (p *Poll) AddColumn(col *Column) (*Column, error) {
	if len(p.Columns) >= 26 {
		return nil, errors.New("columns will exceed the max column count")
	}
	p.Columns = append(p.Columns, col)
	return col, nil
}

func (p *Poll) GetColumn(columnID int) (*Column, error) {
	if columnID < 0 || columnID >= len(p.Columns) {
		return nil, errors.New("invalid column id")
	}
	return p.Columns[columnID], nil
}

func (p *Poll) RemoveColumn(n int) {
	p.Columns = append(p.Columns[:n], p.Columns[n+1:]...)
}

func (p *Poll) AddVote(columnID int) error {
	col, err := p.GetColumn(columnID)
	if err != nil {
		return err
	}
	col.voters++
	return nil
}

func (p *Poll) AddVotes(columnID, n int) error {
	col, err := p.GetColumn(columnID)
	if err != nil {
		return err
	}
	col.voters += n
	return nil
}

func (p *Poll) RemoveVote(columnID int) error {
	col, err := p.GetColumn(columnID)
	if err != nil {
		return err
	}
	col.voters--
	return nil
}

func (p *Poll) RemoveVotes(columnID, n int) error {
	col, err := p.GetColumn(columnID)
	if err != nil {
		return err
	}
	col.voters -= n
	if col.voters < 0 {
		col.voters = 0
	}
	return nil
}

func (p *Poll) GetAllVotes() int {
	c := 0

	for _, col := range p.Columns {
		c += col.VoteCount()
	}
	return c
}

func (p *Poll) GetVotes(columnID int) (int, error) {
	col, err := p.GetColumn(columnID)
	if err != nil {
		return 0, err
	}
	return col.VoteCount(), nil
}

func (p *Poll) ClearVotes() {
	for _, col := range p.Columns {
		col.voters = 0
	}
}
