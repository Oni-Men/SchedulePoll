package bot

import "github.com/bwmarrin/discordgo"

var helpEmbedBuilder *EmbedBuilder

func init() {
	b := NewEmbedBuilder()
	b.Title("ぽる助の使い方")
	b.AddField(&discordgo.MessageEmbedField{
		Name:  "/" + NameHelpCommand,
		Value: "このヘルプを表示します",
	})
	b.AddField(&discordgo.MessageEmbedField{
		Name:  "/" + NamePollCommand,
		Value: "投票の作成、編集、削除ができます。",
	})
	b.AddField(&discordgo.MessageEmbedField{
		Name:  "/" + NameVoteCommand,
		Value: "投票に参加することができます。",
	})
	helpEmbedBuilder = b
}

var HelpCommand = BotCommand{
	NameHelpCommand,
	discordgo.ApplicationCommand{
		Name:        NameHelpCommand,
		Description: "Botの使い方を表示します",
	},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		embed := helpEmbedBuilder.Build(s)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	},
}
