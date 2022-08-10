package printer

import (
	"strconv"

	"github.com/Oni-Men/SchedulePoll/emoji"
	"github.com/Oni-Men/SchedulePoll/poll"
	"github.com/bwmarrin/discordgo"
)

const EMBED_TITLE = ":calendar_spiral: 予定投票 :calendar_spiral:"

func PrintPoll(p *poll.Poll) *discordgo.MessageEmbed {
	var c int
	b := NewEmbedBuilder()
	b.Title(EMBED_TITLE)
	b.Description("#" + p.ID)
	b.FooterText("Discordの都合でグラフへの反映が遅れることがあります。")

	mapping := columnsByYear(p)
	allVotes := p.GetAllVotes()

	for year, list := range mapping {
		var value string
		for _, col := range list {
			emoji := emoji.ABCs[c]
			progress := createProgress(float64(col.VoteCount())/float64(allVotes), 20)
			value += emoji + " **" + col.When.Format("01/02") + "** " + progress + "\n\n"
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
