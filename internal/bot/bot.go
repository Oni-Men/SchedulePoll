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

type InteractionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Bot struct {
	session           *discordgo.Session
	commands          []slashcmd.ISlashCommand
	componentHandlers map[string]InteractionHandler
	commandHandlers   map[string]InteractionHandler
	modalHandlers     map[string]InteractionHandler
	cleanup           bool
}

// 与えられたトークンからBotを生成し、返します
func Create(token string, cleanupCommands bool) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to create a new instance: ", err)
	}

	return &Bot{
		session:           s,
		commands:          []slashcmd.ISlashCommand{},
		componentHandlers: make(map[string]InteractionHandler),
		commandHandlers:   make(map[string]InteractionHandler),
		modalHandlers:     make(map[string]InteractionHandler),
		cleanup:           cleanupCommands,
	}
}

// Botにインテンツを追加します
func (b *Bot) AddIntents(i discordgo.Intent) {
	b.session.Identify.Intents |= i
}

// Botを開始します
func (b *Bot) Start() int {
	guildID := os.Getenv("GUILD_ID")
	b.AddIntents(discordgo.IntentGuildMessages)
	b.AddIntents(discordgo.IntentGuildMessageReactions)

	err := b.session.Open()
	if err != nil {
		fmt.Println("Failed to open a new connection: ", err)
		return 1
	}
	defer b.session.Close()

	registerdCommands := make([]*discordgo.ApplicationCommand, len(b.commands))
	for i, cmd := range b.commands {
		appCmd, err := b.addApplicationCommand(guildID, &cmd)
		if err != nil {
			log.Println(err)
			continue
		}
		registerdCommands[i] = appCmd
	}

	b.session.AddHandler(b.handleInteraction)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	for _, cmd := range registerdCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.Application.ID, guildID, cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", cmd.Name, err)
		}
	}

	return 0
}

func (b *Bot) AddHandler(handler any) {
	b.session.AddHandler(handler)
}

func (b *Bot) AddComponentHandler(id string, handler InteractionHandler) error {
	if _, ok := b.componentHandlers[id]; ok {
		return fmt.Errorf("component handler %s has already existed", id)
	}
	b.componentHandlers[id] = handler

	log.Printf("[Component Handler] %s has been successfully registered.", id)
	return nil
}

func (b *Bot) AddModalHandler(id string, handler InteractionHandler) error {
	if _, ok := b.modalHandlers[id]; ok {
		return fmt.Errorf("modal handler %s has already existed", id)
	}
	b.modalHandlers[id] = handler

	log.Printf("[Modal Handler] %s has been successfully registered.", id)
	return nil
}

func (b *Bot) addApplicationCommand(guildID string, cmd *slashcmd.ISlashCommand) (appCmd *discordgo.ApplicationCommand, err error) {
	id := b.session.State.User.ID
	appCmd = slashcmd.ToApplicationCommand(*cmd)
	appCmd, err = b.session.ApplicationCommandCreate(id, guildID, appCmd)
	if err != nil {
		return
	}

	b.commandHandlers[appCmd.ID] = (*cmd).Handle
	return
}

func (b *Bot) AddSlashCommand(cmd slashcmd.ISlashCommand) {
	b.commands = append(b.commands, cmd)
}

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleSlashCommand := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmdData := i.ApplicationCommandData()
		cmdID := cmdData.ID
		if handler, ok := b.commandHandlers[cmdID]; ok {
			handler(s, i)
		}
	}

	handleMessageComponent := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		msgData := i.MessageComponentData()
		msgID := msgData.CustomID
		if handler, ok := b.componentHandlers[msgID]; ok {
			handler(s, i)
		}
	}

	handleModalSubmit := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		modalData := i.ModalSubmitData()
		modalID := modalData.CustomID
		if handler, ok := b.modalHandlers[modalID]; ok {
			handler(s, i)
		}
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleSlashCommand(s, i)
	case discordgo.InteractionMessageComponent:
		handleMessageComponent(s, i)
	case discordgo.InteractionModalSubmit:
		handleModalSubmit(s, i)
	}
}
