package poll

import (
	"time"
)

type Column struct {
	When   time.Time
	Long   time.Duration
	voters int
}

func CreateColumn(when time.Time, long time.Duration) *Column {
	if long < 0 {
		long = 24 * time.Hour
	}
	col := Column{
		When:   when,
		Long:   long,
		voters: 0,
	}

	return &col
}

func (col *Column) VoteCount() int {
	return col.voters
}
