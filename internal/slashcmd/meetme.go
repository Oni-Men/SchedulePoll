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
			Title:    "MeetMe Form",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "poll-title",
							Label:     "Summary",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 50,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "poll-description",
							Label:     "Description",
							Style:     discordgo.TextInputParagraph,
							Required:  false,
							MaxLength: 500,
							MinLength: 0,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "poll-due",
							Label:       "Due",
							Style:       discordgo.TextInputShort,
							Placeholder: time.Now().AddDate(0, 0, 7).Format("2006/1/2 15:04"),
							Required:    false,
							MaxLength:   16,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "poll-date-list",
							Label:     "Candidates",
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
