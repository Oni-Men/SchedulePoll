package bot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Oni-Men/SchedulePoll/emoji"
	"github.com/Oni-Men/SchedulePoll/poll"
	"github.com/Oni-Men/SchedulePoll/printer"
	"github.com/bwmarrin/discordgo"
)

type BotCommand struct {
	Name    string
	Command discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func Create(botToken string) *discordgo.Session {
	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Failed to create a new instance: ", err)
	}

	return s
}

func Start(s *discordgo.Session) int {
	s.AddHandler(HandleMessageCreate)
	s.AddHandler(HandleReactionAdd)
	s.AddHandler(HandleReactionRemove)
	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	err := s.Open()
	if err != nil {
		fmt.Println("Failed to open a new connection: ", err)
		return 1
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	s.Close()
	return 0
}

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	columns, err := ParseScheduleInput(m.Content)
	if columns == nil {
		return
	}

	if err != nil {
		fmt.Println(err)
		_, err = s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.MessageReference)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	p := poll.CreatePoll()
	p.AddColumnsAll(columns)
	poll.AddPoll(p)

	embed := printer.PrintPoll(p)
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println(err)
	}

	//Add reactions to vote
	for i := 0; i < len(p.Columns); i++ {
		emoji := emoji.ABCs[i]
		err = s.MessageReactionAdd(m.ChannelID, msg.ID, emoji)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func HandleReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := getPollFromMessage(msg)
	if p == nil {
		return
	}

	removeInvalidReactions(p, s, m.ChannelID, m.MessageID, m.Emoji.Name)
	p.AddVote(getIndexFromEmoji(m.Emoji.Name))

	syncVoterCount(p, s, m.ChannelID, m.MessageID)

	embed := printer.PrintPoll(p)
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
	if err != nil {
		fmt.Println(err)
	}
}

func HandleReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.UserID == s.State.User.ID {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := getPollFromMessage(msg)
	if p == nil {
		return
	}

	removeInvalidReactions(p, s, m.ChannelID, m.MessageID, m.Emoji.Name)
	p.RemoveVote(getIndexFromEmoji(m.Emoji.Name))

	syncVoterCount(p, s, m.ChannelID, m.MessageID)

	embed := printer.PrintPoll(p)
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
	if err != nil {
		fmt.Println(err)
	}

}

func getPollFromMessage(msg *discordgo.Message) *poll.Poll {
	if len(msg.Embeds) != 1 {
		return nil
	}

	p, err := getPollFromEmbed(msg.Embeds[0])
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return p
}

func getPollFromEmbed(embed *discordgo.MessageEmbed) (*poll.Poll, error) {
	var err error
	if !strings.HasPrefix(embed.Description, "#") {
		return nil, errors.New("invalid format")
	}

	id := strings.TrimPrefix(embed.Description, "#")
	p := poll.GetPoll(id)

	if p == nil {
		if p, err = ParsePollEmbed(embed); err != nil {
			return nil, err
		}
		poll.AddPoll(p)
	}

	return p, nil
}

func removeInvalidReactions(p *poll.Poll, s *discordgo.Session, channelID, messageID, emojiID string) {
	if !isValidReaction(emojiID, p) {
		err := s.MessageReactionsRemoveEmoji(channelID, messageID, emojiID)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}

func getIndexFromEmoji(id string) int {
	return int([]rune(id)[0] - '\U0001F1E6')
}

func syncVoterCount(p *poll.Poll, s *discordgo.Session, channelID, messageID string) {
	msg, err := s.ChannelMessage(channelID, messageID)
	if err != nil {
		fmt.Println(err)
	}

	p.ClearVotes()

	for _, r := range msg.Reactions {
		if !isValidReaction(r.Emoji.Name, p) {
			continue
		}

		p.AddVotes(getABCIndex(r.Emoji.Name), r.Count-1)
	}
}

func getABCIndex(emojiID string) int {
	if emojiID == "" {
		return -1
	}
	return int([]rune(emojiID)[0]) - int([]rune(emoji.ABCs[0])[0])
}

func isValidReaction(emojiID string, p *poll.Poll) bool {
	n := getABCIndex(emojiID)
	return n >= 0 && n < len(p.Columns)
}
