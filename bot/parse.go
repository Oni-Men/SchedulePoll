package bot

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Oni-Men/SchedulePoll/poll"
	"github.com/Oni-Men/SchedulePoll/printer"
	"github.com/bwmarrin/discordgo"
)

const YOTEI_PREFIX = "!yotei"

func ParseScheduleInput(content string) ([]time.Time, error) {
	var err error

	if !strings.HasPrefix(content, YOTEI_PREFIX) {
		return nil, nil
	}

	content = strings.TrimPrefix(content, YOTEI_PREFIX)

	rawDateList := strings.Split(content, ",")
	parsedDateList := make([]time.Time, 0, len(rawDateList))

	var t time.Time
	year, month, day := time.Now().Date()
	for _, input := range rawDateList {
		split := strings.Split(strings.TrimSpace(input), "/")

		if len(split) == 3 {
			year, err = parseYear(split)
			if err != nil {
				break
			}
			split = split[1:] // Pop "year" element
		}

		if len(split) == 2 {
			month, err = parseMonth(split)
			if err != nil {
				break
			}
			split = split[1:] // Pop "month" element
		}

		if len(split) == 1 {
			if day, err = strconv.Atoi(split[0]); err != nil {
				break
			}
		}

		t = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		parsedDateList = append(parsedDateList, t)
	}

	return parsedDateList, err
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

		date, err := parseEmbedDate(year, field.Value)
		if err != nil {
			return nil, err
		}

		p.AddColumnsAll(date)
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

func parseEmbedDate(year int, text string) ([]time.Time, error) {

	lines := strings.Split(text, "\n")
	parsed := make([]time.Time, 0, len(lines)/2)
	for _, line := range lines {
		if line == "" {
			continue
		}

		// 🇦 **08/01** ◼️◼️◼️ => [🇦, **08/01**, ◼️◼️◼️]
		split := strings.Split(line, " ")
		if len(split) != 3 {
			return nil, errors.New("invalid date format #1")
		}

		// **08/01** => 08/01
		unformatted := strings.Trim(split[1], "*")

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

		date, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}

		t := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.Local)
		parsed = append(parsed, t)
	}

	return parsed, nil
}
