package dateparser

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	var currTime = time.Now()
	var currYear, _, _ = currTime.Date()

	testdata := []struct {
		input  string
		expect []ParsedDateResult
	}{
		{
			`
			5/1,2,3
			`,
			[]ParsedDateResult{
				{
					Date:    time.Date(currYear, 5, 1, 0, 0, 0, 0, time.Local),
					BeginAt: 0,
					EndAt:   24 * time.Hour,
				},
				{
					Date:    time.Date(currYear, 5, 2, 0, 0, 0, 0, time.Local),
					BeginAt: 0,
					EndAt:   24 * time.Hour,
				},
				{
					Date:    time.Date(currYear, 5, 3, 0, 0, 0, 0, time.Local),
					BeginAt: 0,
					EndAt:   24 * time.Hour,
				},
			},
		},
		{
			`
			1/2 12:00-13:00
			1/3 12:00-13:00
			1/4 14:00-15:00
			1/5 14:00-15:00
			`,
			[]ParsedDateResult{
				{
					Date:    time.Date(currYear, 1, 2, 0, 0, 0, 0, time.Local),
					BeginAt: 12 * time.Hour,
					EndAt:   13 * time.Hour,
				},
				{
					Date:    time.Date(currYear, 1, 3, 0, 0, 0, 0, time.Local),
					BeginAt: 12 * time.Hour,
					EndAt:   13 * time.Hour,
				},
				{
					Date:    time.Date(currYear, 1, 4, 0, 0, 0, 0, time.Local),
					BeginAt: 14 * time.Hour,
					EndAt:   15 * time.Hour,
				},
				{
					Date:    time.Date(currYear, 1, 5, 0, 0, 0, 0, time.Local),
					BeginAt: 14 * time.Hour,
					EndAt:   15 * time.Hour,
				},
			},
		},
		{
			`4/1 15:00`,
			[]ParsedDateResult{
				{
					Date:    time.Date(currYear, 4, 1, 0, 0, 0, 0, time.Local),
					BeginAt: 15 * time.Hour,
					EndAt:   24 * time.Hour,
				},
			},
		},
		{
			`
			2023/2/1
			2
			3
			`,
			[]ParsedDateResult{
				{
					Date:    time.Date(2023, 2, 1, 0, 0, 0, 0, time.Local),
					BeginAt: 0,
					EndAt:   24 * time.Hour,
				},
				{
					Date:    time.Date(2023, 2, 2, 0, 0, 0, 0, time.Local),
					BeginAt: 0,
					EndAt:   24 * time.Hour,
				},
				{
					Date:    time.Date(2023, 2, 3, 0, 0, 0, 0, time.Local),
					BeginAt: 0,
					EndAt:   24 * time.Hour,
				},
			},
		},
	}

	for _, td := range testdata {
		parser := NewDateParser(td.input)
		actuals := make([]ParsedDateResult, 0, len(td.expect))

		for parser.HasNext() {
			res, err := parser.Next()

			if err != nil {
				t.Fatal(err)
			}

			if res == nil {
				continue
			}

			actuals = append(actuals, *res)
		}

		testResultEquality(t, actuals, td.expect)

	}
}

func testResultEquality(t *testing.T, actuals, expects []ParsedDateResult) {
	if len(actuals) != len(expects) {
		t.Fatalf("len is not the same. actual=%d, expect=%d", len(actuals), len(expects))
	}

	for i := 0; i < len(expects); i++ {
		actual := actuals[i]
		expect := expects[i]

		if !actual.Date.Equal(expect.Date) {
			layout := "2006/01/02"
			t.Fatalf("when: expects %s but got %s", expect.Date.Format(layout), actual.Date.Format(layout))
		}

		if actual.BeginAt != expect.BeginAt {
			t.Fatalf("begin: expects %s but got %s", expect.BeginAt.String(), actual.BeginAt.String())
		}

		if actual.EndAt != expect.EndAt {
			t.Fatalf("end: expects %s but got %s", expect.EndAt.String(), actual.EndAt.String())
		}
	}
}
