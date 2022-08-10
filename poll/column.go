package poll

import (
	"time"
)

type Column struct {
	When   time.Time
	voters int
}

func CreateColumn(when time.Time) *Column {
	col := Column{
		When:   when,
		voters: 0,
	}

	return &col
}

func (col *Column) VoteCount() int {
	return col.voters
}
