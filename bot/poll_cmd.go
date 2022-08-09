package bot

import "github.com/bwmarrin/discordgo"

var PollCommand = BotCommand{
	NamePollCommand,
	discordgo.ApplicationCommand{
		Name:        NamePollCommand,
		Description: "投票の作成ができます",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "予定(カンマ区切り、年と月は省略可能)",
				Description: "例: 2022/1,2,2/2,3",
				Required:    true,
			},
		},
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	},
}
