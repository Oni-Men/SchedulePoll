package slashcmd

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type YoteiCommand struct{}

var _ ISlashCommand = (*YoteiCommand)(nil)

func (cmd *YoteiCommand) ID() string {
	return "bot.command.yotei"
}

func (cmd *YoteiCommand) Name() string {
	return "yotei"
}

func (cmd *YoteiCommand) Version() string {
	return "1.0.0"
}

func (cmd *YoteiCommand) Description() string {
	return "Create a new Yotei Poll"
}

func (cmd *YoteiCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (cmd *YoteiCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
							MaxLength:   300,
							MinLength:   3,
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
							MaxLength:   1000,
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
							Placeholder: time.Now().Format("2006/01/02 15:04"),
							Required:    false,
							MaxLength:   20,
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
