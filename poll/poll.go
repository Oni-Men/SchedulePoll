package poll

import (
	"errors"

	"github.com/google/uuid"
)

type Voter string

type Poll struct {
	ID       string
	Columns  []*Column
	MaxVotes int
}

func CreatePoll() *Poll {
	p := Poll{
		ID:      uuid.NewString(),
		Columns: make([]*Column, 0, 10),
	}
	return &p
}

func (p *Poll) AddColumn(content string) (*Column, error) {
	if len(p.Columns) >= 10 {
		return nil, errors.New("columns will exceed the max column count")
	}
	col := CreateColumn(content)
	p.Columns = append(p.Columns, col)
	return col, nil
}

func (p *Poll) GetColumn(columnID int) (*Column, error) {
	if columnID < 0 || columnID >= len(p.Columns) {
		return nil, errors.New("invalid column id")
	}
	return p.Columns[columnID], nil
}

func (p *Poll) RemoveColumn(col *Column) {
	for i, c := range p.Columns {
		if c.ID == col.ID {
			p.Columns = append(p.Columns[:i], p.Columns[i+1:]...)
			return
		}
	}
}

func (p *Poll) AddVote(v Voter, columnID int) error {
	col, err := p.GetColumn(columnID)
	if err != nil {
		return err
	}
	col.AppendVoter(v)
	return nil
}

func (p *Poll) RemoveVote(v Voter, columnID int) {

}

func (p *Poll) GetVotes(columnID int) (int, error) {
	col, err := p.GetColumn(columnID)
	if err != nil {
		return 0, err
	}
	return col.VoteCount(), nil
}
