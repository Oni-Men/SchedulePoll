package bot

import (
	"strconv"
	"testing"
	"time"
)

func TestParseSchedule(t *testing.T) {
	y := strconv.Itoa(time.Now().Year())
	testdata := []struct {
		Content string
		When    []string
	}{
		{
			"!予定投票: 8/10,11,14,9/1,2,3",
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
			"!予定投票: 2022/8/10,11,2023/1/2,3",
			[]string{
				"2022/08/10",
				"2022/08/11",
				"2023/01/02",
				"2023/01/03",
			},
		},
	}

	for _, td := range testdata {
		parsed, err := ParseSchedule(td.Content)

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
