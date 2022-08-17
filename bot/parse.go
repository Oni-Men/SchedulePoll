package bot

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Oni-Men/SchedulePoll/poll"
	"github.com/Oni-Men/SchedulePoll/printer"
	"github.com/bwmarrin/discordgo"
)

const YOTEI_PREFIX = "!yotei"

func ParseScheduleInput(content string) ([]*poll.Column, error) {
	var err error
	var startAt, endAt time.Duration
	if !strings.HasPrefix(content, YOTEI_PREFIX) {
		return nil, nil
	}

	content = strings.TrimPrefix(content, YOTEI_PREFIX)
	dateList := strings.Split(content, ",")
	columns := make([]*poll.Column, 0, len(dateList))

	year, month, date := time.Now().Date()
	for _, input := range dateList {
		input = strings.TrimSpace(input)
		split := strings.Split(input, "/")

		if len(split) == 3 {
			year, err = parseYear(split)
			if err != nil {
				break
			}
			split = pop(split) // Pop "year" element
		}

		if len(split) == 2 {
			month, err = parseMonth(split)
			if err != nil {
				break
			}
			split = pop(split) // Pop "month" element
		}

		if len(split) == 1 {
			dateText, durationText := getDateAndDuration(split[0])
			date, err = parseDate(dateText)
			if err != nil {
				break
			}

			startAt, endAt, err = parseStartAndEnd(durationText)
			if err != nil {
				break
			}
		}

		when := time.Date(year, month, date, 0, 0, 0, 0, time.Local)
		when = when.Add(startAt)
		columns = append(columns, poll.CreateColumn(when, endAt-startAt))
	}

	return columns, err
}

func parseYear(split []string) (int, error) {
	year, err := strconv.Atoi(split[0])
	if err != nil {
		return -1, err
	}
	if year < 0 {
		return -1, errors.New("西暦は0以上である必要があります")
	}
	return year, nil
}

func parseMonth(split []string) (time.Month, error) {
	m, err := strconv.Atoi(split[0])
	if err != nil {
		return -1, err
	}
	if m < 1 || m > 12 {
		return -1, errors.New("月は1~12で指定する必要があります")
	}
	return time.Month(m), nil
}

func getDateAndDuration(str string) (string, string) {
	re := regexp.MustCompile(`^(\d{1,2})(?:\[(.+)\])?$`)
	group := re.FindSubmatch([]byte(str))
	if group == nil {
		return "", ""
	}

	return string(group[1]), string(group[2])
}

func parseDate(str string) (int, error) {
	date, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return date, nil
}

func parseStartAndEnd(str string) (time.Duration, time.Duration, error) {
	split := strings.Split(str, "-")
	if len(split) != 2 {
		return 0, 0, nil
	}

	startAt, err := time.Parse("15:04", split[0])
	if err != nil {
		return 0, 0, err
	}

	endAt, err := time.Parse("15:04", split[1])
	if err != nil {
		return 0, 0, err
	}

	zero := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	return startAt.Sub(zero), endAt.Sub(zero), nil
}

func ParsePollEmbed(embed *discordgo.MessageEmbed) (*poll.Poll, error) {

	if embed.Title != printer.EMBED_TITLE {
		return nil, nil
	}

	if !strings.HasPrefix(embed.Description, "#") {
		return nil, nil
	}

	id := strings.TrimPrefix(embed.Description, "#")
	p := poll.CreatePoll()
	p.ID = id

	for _, field := range embed.Fields {
		year, err := parseEmbedYear(field.Name)
		if err != nil {
			return nil, err
		}

		columns, err := parseEmbed(year, field.Value)
		if err != nil {
			return nil, err
		}

		p.AddColumnsAll(columns)
	}

	sort.Slice(p.Columns, func(i, j int) bool {
		return p.Columns[i].When.Before(p.Columns[j].When)
	})

	return p, nil
}

func parseEmbedYear(text string) (int, error) {
	if !strings.HasSuffix(text, "年") {
		return -1, errors.New("invalid year format")
	}

	text = strings.TrimSuffix(text, "年")
	return strconv.Atoi(text)
}

func parseEmbed(year int, text string) ([]*poll.Column, error) {
	lines := strings.Split(text, "\n")
	columns := make([]*poll.Column, 0, len(lines)/2)

	for i := 0; i+1 < len(lines); i += 2 {
		when, err := parseUpperLine(lines[i], year)
		if err != nil {
			return nil, err
		}

		long, err := parseLowerLine(lines[i+1])
		if err != nil {
			return nil, err
		}

		columns = append(columns, poll.CreateColumn(*when, long))
	}

	return columns, nil
}

func parseUpperLine(line string, year int) (*time.Time, error) {
	// 絵文字の部分を取り除く
	// 🇦 **08/01** ◼️◼️◼️ => [🇦, **08/01**, ◼️◼️◼️]
	split := strings.Split(line, " ")
	if len(split) != 3 {
		return nil, errors.New("invalid date format #1")
	}

	// 装飾記号を取り除く
	// **08/01** => 08/01
	unformatted := strings.Trim(split[1], "*")

	// カンマで分割して配列にする
	// 08/01 => [08, 01]
	split = strings.Split(unformatted, "/")
	if len(split) != 2 {
		return nil, errors.New("invalit date format #2")
	}

	month, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, err
	}
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}

	date, err := parsePollDate(split[1])
	if err != nil {
		return nil, err
	}

	t := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.Local)
	return &t, nil
}

func parseLowerLine(line string) (time.Duration, error) {
	// 空白を取り除く
	line = strings.Trim(line, " ")
	split := strings.Split(line, "-")
	if len(split) != 2 {
		return 0, errors.New("invalid format: parseLowerLine#1")
	}

	startAt, err := time.Parse("15:04", strings.Trim(split[0], " "))
	if err != nil {
		return 0, err
	}
	endAt, err := time.Parse("15:04", strings.Trim(split[1], " "))
	if err != nil {
		return 0, err
	}

	return endAt.Sub(startAt), nil
}

func parsePollDate(str string) (int, error) {
	re := regexp.MustCompile(`^(\d{2})\(.\)$`)
	group := re.FindSubmatch([]byte(str))
	if group == nil {
		return 0, errors.New("invalid date format")
	}

	date, err := strconv.Atoi(string(group[1]))
	if err != nil {
		return 0, err
	}

	return date, nil
}

func pop(s []string) []string {
	return s[1:]
}
