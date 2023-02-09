package slashcmd

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MeetmeCommand struct{}

var _ ISlashCommand = (*MeetmeCommand)(nil)

func (cmd *MeetmeCommand) ID() string {
	return "bot.command.meetme"
}

func (cmd *MeetmeCommand) Name() string {
	return "meetme"
}

func (cmd *MeetmeCommand) Version() string {
	return "1.0.0"
}

func (cmd *MeetmeCommand) Description() string {
	return "Create a new Poll to find a good time to meet"
}

func (cmd *MeetmeCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (cmd *MeetmeCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: PollCreateModal,
			Title:    "日程投票の作成",
			Components: []discordgo.MessageComponent{
				// 投票のタイトル
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "poll-title",
							Label:       "この投票のタイトル",
							Style:       discordgo.TextInputShort,
							Placeholder: "カラオケ行ける日募集～",
							Required:    true,
							MaxLength:   50,
						},
					},
				},
				// 投票の説明
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "poll-description",
							Label:       "日程の目的など、投票者に伝えたいことを入力してください。",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "集合場所は駅前でーす",
							Required:    false,
							MaxLength:   500,
							MinLength:   0,
						},
					},
				},
				// 投票の期限
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "poll-due",
							Label:       "投票の期限",
							Style:       discordgo.TextInputShort,
							Placeholder: time.Now().AddDate(0, 0, 7).Format("2006/1/2 15:04"),
							Required:    false,
							MaxLength:   16,
						},
					},
				},
				// 候補の日程リスト
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "poll-date-list",
							Label:     "日程リスト",
							Style:     discordgo.TextInputParagraph,
							Required:  true,
							MaxLength: 500,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
	}
}
