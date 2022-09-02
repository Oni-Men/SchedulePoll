package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oni-Men/SchedulePoll/internal/slashcmd"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session  *discordgo.Session
	commands []slashcmd.ISlashCommand
}

func Create(token string) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to create a new instance: ", err)
	}
	return &Bot{
		session:  s,
		commands: []slashcmd.ISlashCommand{},
	}
}

func (b *Bot) AddIntents(i discordgo.Intent) {
	b.session.Identify.Intents |= i
}

func (b *Bot) Start() int {
	b.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	err := b.session.Open()
	if err != nil {
		fmt.Println("Failed to open a new connection: ", err)
		return 1
	}

	guildID := os.Getenv("GUILD_ID")
	for _, cmd := range b.commands {
		b.addApplicationCommand(guildID, &cmd)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	b.session.Close()
	return 0
}

func (b *Bot) AddHandler(handler any) {
	b.session.AddHandler(handler)
}

func (b *Bot) addApplicationCommand(guildID string, cmd *slashcmd.ISlashCommand) {
	id := b.session.State.User.ID
	appCmd := slashcmd.ToApplicationCommand(*cmd)
	_, err := b.session.ApplicationCommandCreate(id, guildID, appCmd)
	if err != nil {
		fmt.Println(err)
	}

	b.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		(*cmd).Handle(s, i)
	})
}

func (b *Bot) AddSlashCommand(cmd slashcmd.ISlashCommand) {
	b.commands = append(b.commands, cmd)
}

func (b *Bot) RespondText(i *discordgo.Interaction, text string) {
	b.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	})
}
