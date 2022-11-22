package printer

import "github.com/bwmarrin/discordgo"

// A builder for discord-embed
type EmbedPrinter struct {
	emb *discordgo.MessageEmbed
}

func New() *EmbedPrinter {
	return &EmbedPrinter{&discordgo.MessageEmbed{}}
}

func (e *EmbedPrinter) Title(v string) *EmbedPrinter {
	e.emb.Title = v
	return e
}

func (e *EmbedPrinter) Description(v string) *EmbedPrinter {
	e.emb.Description = v
	return e
}

func (e *EmbedPrinter) FooterText(text string) *EmbedPrinter {
	if e.emb.Footer == nil {
		e.emb.Footer = &discordgo.MessageEmbedFooter{}
	}
	e.emb.Footer.Text = text
	return e
}

func (e *EmbedPrinter) FooterIcon(url string) *EmbedPrinter {
	if e.emb.Footer == nil {
		e.emb.Footer = &discordgo.MessageEmbedFooter{}
	}
	e.emb.Footer.IconURL = url
	return e
}

func (e *EmbedPrinter) AddField(name, value string, inline ...bool) *EmbedPrinter {
	if value == "" {
		value = "---"
	}
	field := &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: len(inline) > 0 && inline[0],
	}
	if e.emb.Fields == nil {
		e.emb.Fields = make([]*discordgo.MessageEmbedField, 0, 10)
	}
	e.emb.Fields = append(e.emb.Fields, field)
	return e
}

func (e *EmbedPrinter) AddInlineField(name, value string) *EmbedPrinter {
	e.AddField(name, name, true)
	return e
}

func (e *EmbedPrinter) Build() *discordgo.MessageEmbed {
	return e.emb
}
