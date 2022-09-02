package poll

import (
	"strconv"

	"github.com/Oni-Men/SchedulePoll/pkg/emoji"
	"github.com/Oni-Men/SchedulePoll/pkg/printer"
	"github.com/Oni-Men/SchedulePoll/pkg/timeutil"
	"github.com/bwmarrin/discordgo"
)

const EMBED_TITLE = ":calendar_spiral: 予定投票 :calendar_spiral:"

var WeekDays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func PrintPoll(p *Poll) *discordgo.MessageEmbed {
	var c int
	b := printer.New()
	b.Title(EMBED_TITLE)
	b.Description("#" + p.ID)
	b.FooterText("Discordの都合でグラフへの反映が遅れることがあります。")

	mapping := columnsByYear(p)
	allVotes := float64(p.GetAllVotes())

	for year, columns := range mapping {
		value := ""
		for _, col := range columns {
			emoji := emoji.ABCDEmoji(c)
			ratio := float64(col.VoteCount()) / allVotes
			progress := createProgress(ratio, 20)
			weekDayText := WeekDays[col.Date.Weekday()]
			dateText := col.Date.Format("01/02") + "(" + weekDayText + ")"
			beginAt := timeutil.GetZeroTime().Add(col.BeginAt).Format("15:04")
			endAt := timeutil.GetZeroTime().Add(col.EndAt).Format("15:04")

			value += emoji + " **" + dateText + "** " + progress + "\n"
			value += "    " + beginAt + " - " + endAt + "\n"
			c++
		}

		b.AddField(strconv.Itoa(year)+"年", value)
	}

	return b.Build()
}

func createProgress(ratio float64, length int) string {
	var text string
	n := int(ratio * float64(length))
	for i := 0; i < length; i++ {
		if i < n {
			text += emoji.ProgressFG
		} else {
			text += emoji.ProgressBG
		}
	}
	return text
}

func columnsByYear(p *Poll) map[int][]*Column {
	mapping := make(map[int][]*Column)

	for _, col := range p.Columns {
		list := mapping[col.Date.Year()]
		if list == nil {
			list = make([]*Column, 0, 5)
		}
		mapping[col.Date.Year()] = append(list, col)
	}

	return mapping
}
