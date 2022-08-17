package printer

import (
	"strconv"

	"github.com/Oni-Men/SchedulePoll/emoji"
	"github.com/Oni-Men/SchedulePoll/poll"
	"github.com/bwmarrin/discordgo"
)

const EMBED_TITLE = ":calendar_spiral: 予定投票 :calendar_spiral:"

var WeekDays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func PrintPoll(p *poll.Poll) *discordgo.MessageEmbed {
	var c int
	b := NewEmbedBuilder()
	b.Title(EMBED_TITLE)
	b.Description("#" + p.ID)
	b.FooterText("Discordの都合でグラフへの反映が遅れることがあります。")

	mapping := columnsByYear(p)
	allVotes := float64(p.GetAllVotes())

	for year, columns := range mapping {
		value := ""
		for _, col := range columns {
			emoji := emoji.ABCs[c]
			ratio := float64(col.VoteCount()) / allVotes
			progress := createProgress(ratio, 20)
			weekDayText := WeekDays[col.When.Weekday()]
			dateText := col.When.Format("01/02") + "(" + weekDayText + ")"
			startAt := col.When.Format("15:04")
			endAt := col.When.Add(col.Long).Format("15:04")

			value += emoji + " **" + dateText + "** " + progress + "\n"
			value += "    " + startAt + " - " + endAt + "\n"
			c++
		}

		b.AddField(&discordgo.MessageEmbedField{
			Name:   strconv.Itoa(year) + "年",
			Value:  value,
			Inline: false,
		})
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

func columnsByYear(p *poll.Poll) map[int][]*poll.Column {
	mapping := make(map[int][]*poll.Column)

	for _, col := range p.Columns {
		list := mapping[col.When.Year()]
		if list == nil {
			list = make([]*poll.Column, 0, 5)
		}
		mapping[col.When.Year()] = append(list, col)
	}

	return mapping
}
