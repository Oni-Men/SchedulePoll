package bot

import "github.com/bwmarrin/discordgo"

type EmbedBuilder struct {
	title       string
	description string
	fields      []*discordgo.MessageEmbedField
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

func (e *EmbedBuilder) AddField(v *discordgo.MessageEmbedField) *EmbedBuilder {
	e.fields = append(e.fields, v)
	return e
}

func (e *EmbedBuilder) Build(s *discordgo.Session) *discordgo.MessageEmbed {
	res := discordgo.MessageEmbed{
		Title:       e.title,
		Description: e.description,
		Fields:      e.fields,
	}

	return &res
}
