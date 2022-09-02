package poll

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Oni-Men/SchedulePoll/pkg/timeutil"
	"github.com/bwmarrin/discordgo"
)

func ParsePollEmbed(embed *discordgo.MessageEmbed) (*Poll, error) {
	if embed.Title != EMBED_TITLE {
		return nil, nil
	}

	if !strings.HasPrefix(embed.Description, "#") {
		return nil, nil
	}

	id := strings.TrimPrefix(embed.Description, "#")
	p := CreatePoll()
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

	//日付を古い順でソート
	sort.Slice(p.Columns, func(i, j int) bool {
		return p.Columns[i].Date.Before(p.Columns[j].Date)
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

func parseEmbed(year int, text string) ([]*Column, error) {
	lines := strings.Split(text, "\n")
	columns := make([]*Column, 0, len(lines)/2)

	for i := 0; i+1 < len(lines); i += 2 {
		date, err := parseUpperLine(lines[i], year)
		if err != nil {
			return nil, err
		}

		beginAt, endAt, err := parseLowerLine(lines[i+1])
		if err != nil {
			return nil, err
		}

		columns = append(columns, CreateColumn(*date, beginAt, endAt))
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

func parseLowerLine(line string) (time.Duration, time.Duration, error) {
	// 空白を取り除く
	line = strings.Trim(line, " ")
	split := strings.Split(line, "-")
	if len(split) != 2 {
		return 0, 0, errors.New("invalid format: parseLowerLine#1")
	}

	beginAt, err := time.Parse("15:04", strings.Trim(split[0], " "))
	if err != nil {
		return 0, 0, err
	}
	endAt, err := time.Parse("15:04", strings.Trim(split[1], " "))
	if err != nil {
		return 0, 0, err
	}

	return timeutil.GetElapsedFromZero(beginAt), timeutil.GetElapsedFromZero(endAt), nil
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
