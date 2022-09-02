package dateparser

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Oni-Men/SchedulePoll/pkg/sliceutil"
	"github.com/Oni-Men/SchedulePoll/pkg/timeutil"
)

const (
	YoteiDateFormat = "2006/01/02"
	DurationFormat  = "15:04"
)

// DateParser parses date list.
// Each date in the list is splited by new line
type DateParser struct {
	list      []string
	currYear  int
	currMonth time.Month
	currDay   int

	currIdx int
}

type ParsedDateResult struct {
	Date    time.Time     // This field doesn't have time information.
	BeginAt time.Duration // Elapsed time since 0:00 am
	EndAt   time.Duration // Elapsed time since 0:00 am
}

func NewDateParser(input string) *DateParser {
	var curr = time.Now()
	return &DateParser{
		strings.Split(input, "\n"),
		curr.Year(),
		curr.Month(),
		curr.Day(),
		0,
	}
}

func newParsedDateResult(p *DateParser) *ParsedDateResult {
	return &ParsedDateResult{
		Date:    time.Date(p.currYear, p.currMonth, p.currDay, 0, 0, 0, 0, time.Local),
		BeginAt: 0,
		EndAt:   0,
	}
}

func (p *DateParser) HasNext() bool {
	return p.currIdx < len(p.list)
}

func (p *DateParser) Next() (*ParsedDateResult, error) {
	if p.currIdx >= len(p.list) {
		return nil, errors.New("invalid operation: parser has no more item to parse")
	}

	input := sanitize(p.list[p.currIdx])
	p.currIdx++

	if input == "" {
		return nil, nil
	}

	split := strings.Split(input, "/")

	res := newParsedDateResult(p)

	var pop string
	if len(split) == 3 {
		pop, split = sliceutil.Pop(split) // Pop "year" element
		if err := processYear(res, pop); err != nil {
			return nil, err
		}
	}

	if len(split) == 2 {
		pop, split = sliceutil.Pop(split) // Pop "month" element
		if err := processMonth(res, pop); err != nil {
			return nil, err
		}
	}

	if len(split) == 1 {
		dateStr, durStr := getDateAndDuration(split[0])
		if err := processDate(res, dateStr); err != nil {
			return nil, err
		}

		if err := processDuration(res, durStr); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func processYear(res *ParsedDateResult, str string) error {
	year, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	if year < 0 {
		return errors.New("year must be greater than zero")
	}
	res.Date = timeutil.SetYear(res.Date, year)
	return nil
}

func processMonth(res *ParsedDateResult, str string) error {
	month, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	if month < 1 || month > 12 {
		return errors.New("month must be specified within 1 to 12")
	}
	res.Date = timeutil.SetMonth(res.Date, month)
	return nil
}

func getDateAndDuration(str string) (string, string) {
	split := strings.Split(str, " ")
	if len(split) == 0 {
		return "", ""
	}
	if len(split) == 1 {
		return split[0], ""
	}
	return split[0], split[1]
}

func processDate(res *ParsedDateResult, str string) error {
	day, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	res.Date = timeutil.SetDay(res.Date, day)
	return nil
}

func processDuration(res *ParsedDateResult, str string) error {
	if str == "" {
		return errors.New("invalid duration format")
	}
	split := strings.Split(str, "-")
	if len(split) != 2 {
		return errors.New("invalid duration format")
	}

	if begin, err := time.Parse(DurationFormat, split[0]); err != nil {
		return err
	} else {
		res.BeginAt = timeutil.GetElapsedFromZero(begin)
	}

	if end, err := time.Parse(DurationFormat, split[1]); err != nil {
		return err
	} else {
		res.EndAt = timeutil.GetElapsedFromZero(end)
	}

	return nil
}

func sanitize(str string) string {
	return strings.TrimSpace(str)
}
