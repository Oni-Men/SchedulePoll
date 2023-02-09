package modal

import (
	"time"

	"github.com/Oni-Men/SchedulePoll/pkg/dateparser"
	"github.com/Oni-Men/SchedulePoll/pkg/timeutil"
)

// A struct which represents modal submitted data.
type ModalResponse struct {
	title       string `validate:"max=50"`
	description string `validate:"max=500"`
	due         time.Time
	dateList    []dateparser.ParsedDateResult
}

type ModalResponseOption func(*ModalResponse)

func WithTitle(title string) ModalResponseOption {
	return func(r *ModalResponse) {
		r.title = title
	}
}

func WithDesc(description string) ModalResponseOption {
	return func(r *ModalResponse) {
		r.description = description
	}
}

func WithDue(due time.Time) ModalResponseOption {
	return func(r *ModalResponse) {
		r.due = due
	}
}

func WithDateList(dateList []dateparser.ParsedDateResult) ModalResponseOption {
	return func(r *ModalResponse) {
		r.dateList = dateList
	}
}

func NewModalResponse(options ...ModalResponseOption) *ModalResponse {
	r := &ModalResponse{
		title:       "No title",
		description: "",
		due:         timeutil.GetZeroTime(),
		dateList:    []dateparser.ParsedDateResult{},
	}

	// Applying options
	for _, opt := range options {
		opt(r)
	}

	return r
}

func (r *ModalResponse) Title() string {
	return r.title
}

func (r *ModalResponse) Description() string {
	return r.description
}

func (r *ModalResponse) Due() time.Time {
	return r.due
}

func (r *ModalResponse) DateList() []dateparser.ParsedDateResult {
	return r.dateList
}
