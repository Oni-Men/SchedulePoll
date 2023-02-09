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

// 与えられたトークンからBotを生成し、返します
func Create(token string, cleanupCommands bool) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to create a new instance: ", err)
	}

	return &Bot{
		session:  s,
		commands: []slashcmd.ISlashCommand{},
	}
}

// Botにインテンツを追加します
func (b *Bot) AddIntents(i discordgo.Intent) {
	b.session.Identify.Intents |= i
}

// Botを開始します
func (b *Bot) Start() int {
	b.AddIntents(discordgo.IntentGuildMessages)
	b.AddIntents(discordgo.IntentGuildMessageReactions)

	err := b.session.Open()
	if err != nil {
		fmt.Println("Failed to open a new connection: ", err)
		return 1
	}
	defer b.session.Close()

	guildID := os.Getenv("GUILD_ID")

	commandIds := make(map[string]string, len(b.commands))
	for _, cmd := range b.commands {
		err := b.addApplicationCommand(guildID, &cmd)
		if err != nil {
			log.Println(err)
			continue
		}
		commandIds[cmd.ID()] = cmd.Name()
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	for id, name := range commandIds {
		err := b.session.ApplicationCommandDelete(b.session.State.Application.ID, guildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

	return 0
}

func (b *Bot) AddHandler(handler any) {
	b.session.AddHandler(handler)
}

func (b *Bot) addApplicationCommand(guildID string, cmd *slashcmd.ISlashCommand) error {
	id := b.session.State.User.ID
	appCmd := slashcmd.ToApplicationCommand(*cmd)
	_, err := b.session.ApplicationCommandCreate(id, guildID, appCmd)
	if err != nil {
		return err
	}

	b.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		(*cmd).Handle(s, i)
	})

	return nil
}

func (b *Bot) AddSlashCommand(cmd slashcmd.ISlashCommand) {
	b.commands = append(b.commands, cmd)
}
