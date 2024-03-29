package dateparser

import (
	"errors"
	"regexp"
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

var DateSplitRegexp *regexp.Regexp

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

	// 日付を改行またはコンマで分割する正規表現
	// キャッシュすることで使いまわす
	if DateSplitRegexp == nil {
		var err error
		DateSplitRegexp, err = regexp.Compile("[\n,]")
		if err != nil {
			return nil
		}
	}

	list := DateSplitRegexp.Split(input, -1)
	list = sliceutil.Map(list, func(t string) string {
		return strings.TrimSpace(t)
	})
	list = sliceutil.Filter(list, func(t string) bool {
		return t != ""
	})

	return &DateParser{
		list,
		curr.Year(),
		curr.Month(),
		curr.Day(),
		0,
	}
}

func ParseInlineDate(input string) (*ParsedDateResult, error) {
	dp := NewDateParser(input)
	if dp.HasNext() {
		return dp.Next()
	}
	return nil, errors.New("invalid input value")
}

func newParsedDateResult(p *DateParser) *ParsedDateResult {
	return &ParsedDateResult{
		Date:    time.Date(p.currYear, p.currMonth, p.currDay, 0, 0, 0, 0, time.Local),
		BeginAt: 0,
		EndAt:   24 * time.Hour,
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

	var pop *string
	if len(split) == 3 {
		pop, split, _ = sliceutil.Pop(split) // Pop "year" element
		if year, err := processYear(*pop); err == nil {
			if year < 0 {
				return nil, errors.New("year must be greater than zero")
			}
			res.Date = timeutil.SetYear(res.Date, year)
			p.currYear = year
		} else {
			return nil, err
		}
	}

	if len(split) == 2 {
		pop, split, _ = sliceutil.Pop(split) // Pop "month" element
		if month, err := processMonth(*pop); err == nil {
			if month < 1 || month > 12 {
				return nil, errors.New("month must be specified within 1 to 12")
			}
			res.Date = timeutil.SetMonth(res.Date, month)
			p.currMonth = time.Month(month)
		} else {
			return nil, err
		}
	}

	if len(split) == 1 {
		dateStr, durStr := getDateAndDuration(split[0])
		if date, err := processDate(dateStr); err == nil {
			res.Date = timeutil.SetDay(res.Date, date)
			p.currDay = date
		} else {
			return nil, err
		}

		if err := processDuration(res, durStr); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func processYear(str string) (int, error) {
	year, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}
	return year, nil
}

func processMonth(str string) (int, error) {
	month, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}
	return month, nil
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

func processDate(str string) (int, error) {
	day, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}
	return day, nil
}

func processDuration(res *ParsedDateResult, str string) error {
	if str == "" {
		return nil
	}

	var unpush *string
	split := strings.Split(str, "-")
	if len(split) == 2 {
		unpush, split, _ = sliceutil.Unpush(split)
		if end, err := time.Parse(DurationFormat, *unpush); err != nil {
			return err
		} else {
			res.EndAt = timeutil.GetElapsedFromZero(end)
		}
	}

	if len(split) == 1 {
		unpush, _, _ = sliceutil.Unpush(split)
		if begin, err := time.Parse(DurationFormat, *unpush); err != nil {
			return err
		} else {
			res.BeginAt = timeutil.GetElapsedFromZero(begin)
		}
	}

	return nil
}

func sanitize(str string) string {
	return strings.TrimSpace(str)
}
