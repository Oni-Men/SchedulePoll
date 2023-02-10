package poll

import (
	"time"
)

type Column struct {
	Date    time.Time
	BeginAt time.Duration
	EndAt   time.Duration
	voters  int
}

func CreateColumn(date time.Time, begin, end time.Duration) *Column {
	if end-begin <= 0 {
		begin = 0
		end = 24 * time.Hour
	}
	col := Column{
		Date:    date,
		BeginAt: begin,
		EndAt:   end,
		voters:  0,
	}

	return &col
}

func (col *Column) VoteCount() int {
	return col.voters
}
