package printer

import "github.com/bwmarrin/discordgo"

type EmbedBuilder struct {
	title       string
	description string
	fields      []*discordgo.MessageEmbedField
	footerText  string
	footerIcon  string
}

func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{}
}

func (e *EmbedBuilder) Title(v string) *EmbedBuilder {
	e.title = v
	return e
}

func (e *EmbedBuilder) Description(v string) *EmbedBuilder {
	e.description = v
	return e
}

func (e *EmbedBuilder) FooterText(text string) *EmbedBuilder {
	e.footerText = text
	return e
}

func (e *EmbedBuilder) FooterIcon(url string) *EmbedBuilder {
	e.footerIcon = url
	return e
}

func (e *EmbedBuilder) AddField(v *discordgo.MessageEmbedField) *EmbedBuilder {
	e.fields = append(e.fields, v)
	return e
}

func (e *EmbedBuilder) Build() *discordgo.MessageEmbed {
	res := discordgo.MessageEmbed{
		Title:       e.title,
		Description: e.description,
		Fields:      e.fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    e.footerText,
			IconURL: e.footerIcon,
		},
	}

	return &res
}
