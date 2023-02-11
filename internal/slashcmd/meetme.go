package slashcmd

import (
	"fmt"
	"log"

	"github.com/Oni-Men/SchedulePoll/internal/model"
	"github.com/Oni-Men/SchedulePoll/pkg/emoji"
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
	if i.Member == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can't use this command here",
			},
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("%s, Click the button below to create new poll", i.Member.Mention()),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: model.PollCreateButton,
							Label:    "New poll",
							Style:    discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: emoji.Calendar,
							},
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
