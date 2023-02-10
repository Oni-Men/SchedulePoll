package timeutil

import "time"

func SetYear(t time.Time, year int) time.Time {
	return t.AddDate(year-t.Year(), 0, 0)
}

func SetMonth(t time.Time, month int) time.Time {
	return t.AddDate(0, month-int(t.Month()), 0)
}

func SetDay(t time.Time, day int) time.Time {
	return t.AddDate(0, 0, day-t.Day())
}

func SetHourAndMinutes(t time.Time, hours, minutes int) time.Time {
	h := hours - t.Hour()
	m := minutes - t.Minute()
	return t.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
}

// Returns elapsed duration since zero time instant(0000/01/01 00:00:00)
func GetElapsedFromZero(t time.Time) time.Duration {
	return t.Sub(GetZeroTime())
}

// Returns zero time (0000/01/01 00:00:00).
func GetZeroTime() time.Time {
	return time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
}
