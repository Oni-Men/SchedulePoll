package bot

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/Oni-Men/SchedulePoll/poll"
	"github.com/Oni-Men/SchedulePoll/printer"
)

func TestParseSchedule(t *testing.T) {
	y := strconv.Itoa(time.Now().Year())
	testdata := []struct {
		Case          string
		ExpectedWhens []string
		ExpectedLongs []string
	}{
		{
			"!yotei 8/10,11,14,9/1,2,3",
			[]string{
				y + "/08/10 00:00",
				y + "/08/11 00:00",
				y + "/08/14 00:00",
				y + "/09/01 00:00",
				y + "/09/02 00:00",
				y + "/09/03 00:00",
			},
			[]string{},
		},
		{
			"!yotei 2022/8/10,11,2023/1/2,3",
			[]string{
				"2022/08/10 00:00",
				"2022/08/11 00:00",
				"2023/01/02 00:00",
				"2023/01/03 00:00",
			},
			[]string{},
		},
		{
			"!yotei 8/15[14:30-15:30], 16, 17",
			[]string{
				"2022/08/15 14:30",
				"2022/08/16 00:00",
				"2022/08/17 00:00",
			},
			[]string{
				"1h",
				"0",
				"0",
			},
		},
	}

	for _, td := range testdata {
		parsed, err := ParseScheduleInput(td.Case)

		if err != nil {
			t.Fatalf("no error expected, but we got %v", err)
		}

		if len(parsed) != len(td.ExpectedWhens) {
			t.Fatalf("The length is not the same! %d, %d", len(parsed), len(td.ExpectedWhens))
		}

		for i, d := range parsed {
			formatted := d.When.Format("2006/01/02 15:04")
			if formatted != td.ExpectedWhens[i] {
				t.Fatalf("WRONG!, actual = %s, expected = %s", formatted, td.ExpectedWhens[i])
			}
		}
	}

}

func TestParseEmbed(t *testing.T) {
	testdata := []struct {
		times []time.Time
		longs []time.Duration
	}{
		{
			times: []time.Time{
				time.Date(2022, 8, 10, 0, 0, 0, 0, time.Local),
				time.Date(2022, 8, 11, 0, 0, 0, 0, time.Local),
				time.Date(2022, 9, 12, 0, 0, 0, 0, time.Local),
			},
			longs: []time.Duration{
				0,
				1 * time.Hour,
				2 * time.Minute,
			},
		},
	}

	for _, td := range testdata {
		p := poll.CreatePoll()
		for i, t := range td.times {
			p.AddColumn(poll.CreateColumn(t, td.longs[i]))
		}

		embed := printer.PrintPoll(p)
		q, err := ParsePollEmbed(embed)
		if err != nil {
			log.Fatalf("no error expected but we got %v\n", err)
		}

		if p.ID != q.ID {
			t.Fatalf("ID is not the same. %s, %s", p.ID, q.ID)
		}

		if len(p.Columns) != len(q.Columns) {
			t.Fatalf("The length of columns is not the same. %d, %d", len(p.Columns), len(q.Columns))
		}

		if !isColumnsEqual(p, q) {
			t.Fatalf("Columns are not equal.")
		}
	}
}

func isColumnsEqual(a, b *poll.Poll) bool {
	for i, colA := range a.Columns {
		colB := b.Columns[i]
		if colA.When.Equal(colB.When) {
			continue
		}

		return false
	}

	return true
}

func TestReactionValidator(t *testing.T) {
	p := poll.CreatePoll()
	p.AddColumn(poll.CreateColumn(time.Now(), -1))
	p.AddColumn(poll.CreateColumn(time.Now(), -1))
	p.AddColumn(poll.CreateColumn(time.Now(), -1))

	testdata := []struct {
		id       string
		expected bool
	}{
		{
			"\U0001F1E5",
			false,
		}, {
			"🇦",
			true,
		}, {
			"🇨",
			true,
		}, {
			"🇩",
			false,
		},
	}

	for _, td := range testdata {
		if isValidReaction(td.id, p) != td.expected {
			t.Fatalf("failed to validate for %s.", td.id)
		}
	}
}
