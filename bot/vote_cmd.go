package bot

import "github.com/bwmarrin/discordgo"

var VoteCommand = BotCommand{
	NameVoteCommand,
	discordgo.ApplicationCommand{
		Name:        NameVoteCommand,
		Description: "投票に参加します",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {

	},
}
