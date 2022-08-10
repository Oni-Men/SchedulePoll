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
		Content string
		When    []string
	}{
		{
			"!yotei 8/10,11,14,9/1,2,3",
			[]string{
				y + "/08/10",
				y + "/08/11",
				y + "/08/14",
				y + "/09/01",
				y + "/09/02",
				y + "/09/03",
			},
		},
		{
			"!yotei 2022/8/10,11,2023/1/2,3",
			[]string{
				"2022/08/10",
				"2022/08/11",
				"2023/01/02",
				"2023/01/03",
			},
		},
	}

	for _, td := range testdata {
		parsed, err := ParseScheduleInput(td.Content)

		if err != nil {
			t.Fatalf("no error expected, but we got %v", err)
		}

		if len(parsed) != len(td.When) {
			t.Fatalf("The length is not the same! %d, %d", len(parsed), len(td.When))
		}

		for i, d := range parsed {
			formatted := d.Format("2006/01/02")
			if formatted != td.When[i] {
				t.Fatalf("WRONG!, actual = %s, expected = %s", formatted, td.When[i])
			}
		}
	}

}

func TestParseEmbed(t *testing.T) {
	testdata := []struct {
		times []time.Time
	}{
		{
			times: []time.Time{
				time.Date(2022, 8, 10, 0, 0, 0, 0, time.Local),
				time.Date(2022, 8, 11, 0, 0, 0, 0, time.Local),
				time.Date(2022, 9, 12, 0, 0, 0, 0, time.Local),
			},
		},
	}

	for _, td := range testdata {
		p := poll.CreatePoll()
		p.AddColumnsAll(td.times)
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
	p.AddColumn(time.Now())
	p.AddColumn(time.Now())
	p.AddColumn(time.Now())

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
