package service

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Oni-Men/SchedulePoll/internal/model/poll"
	"github.com/Oni-Men/SchedulePoll/pkg/emoji"
	"github.com/Oni-Men/SchedulePoll/pkg/printer"
	"github.com/Oni-Men/SchedulePoll/pkg/timeutil"
	"github.com/bwmarrin/discordgo"
)

const EMBED_TITLE = ":calendar_spiral: 予定投票 :calendar_spiral:"

var WeekDays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

func PrintPoll(p *poll.Poll) *discordgo.MessageEmbed {
	b := printer.New()
	b.Title(EMBED_TITLE)
	b.Description("#" + p.ID + "\n")
	b.AddField(p.Title, p.Description)

	if p.Due.IsZero() {
		b.AddField("投票期限", "なし")
	} else {
		b.AddField("投票期限", p.Due.Format("2006/01/02 15:04"))
	}

	columnsMap := mapColumnsByYear(p)
	votes := float64(p.GetAllVotes())

	var c int
	for y, columns := range columnsMap {
		var text string
		for _, col := range columns {
			icon := emoji.ABCDEmoji(c)
			ratio := float64(col.VoteCount()) / votes
			progress := createProgress(ratio, 20)
			weekDayText := WeekDays[col.Date.Weekday()]
			ratioText := ratioText(ratio)
			dateText := col.Date.Format("01/02") + "(" + weekDayText + ")"
			beginAt := timeutil.GetZeroTime().Add(col.BeginAt).Format("15:04")
			endAt := timeutil.GetZeroTime().Add(col.EndAt).Format("15:04")

			text += icon + " **" + dateText + "** " + progress + "\n"
			text += "    " + beginAt + " - " + endAt + " " + ratioText + "\n"
			c++
		}

		b.AddField(strconv.Itoa(y)+"年", text)
	}

	b.FooterText("リクエスト制限の都合上グラフへの反映が遅れることがあります。")
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

func mapColumnsByYear(p *poll.Poll) map[int][]*poll.Column {
	mapping := make(map[int][]*poll.Column)

	for _, col := range p.Columns {
		list := mapping[col.Date.Year()]
		if list == nil {
			list = make([]*poll.Column, 0, 5)
		}
		mapping[col.Date.Year()] = append(list, col)
	}

	return mapping
}

func ratioText(ratio float64) string {
	if math.IsNaN(ratio) {
		ratio = 0
	}
	s := strconv.Itoa(int(ratio * 100))
	return fmt.Sprintf(" %3s%% ", s)
}
