package slashcmd

import "github.com/bwmarrin/discordgo"

type ISlashCommand interface {
	ID() string
	Name() string
	Description() string
	Version() string
	Options() []*discordgo.ApplicationCommandOption
	Handle(*discordgo.Session, *discordgo.InteractionCreate)
}

func ToApplicationCommand(cmd ISlashCommand) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		ID:          cmd.ID(),
		Name:        cmd.Name(),
		Version:     cmd.Version(),
		Type:        discordgo.ChatApplicationCommand,
		Description: cmd.Description(),
		Options:     cmd.Options(),
	}
}
