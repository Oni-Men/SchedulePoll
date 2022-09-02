package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Oni-Men/SchedulePoll/internal/bot"
	"github.com/Oni-Men/SchedulePoll/internal/model/modal"
	"github.com/Oni-Men/SchedulePoll/internal/model/poll"
	"github.com/Oni-Men/SchedulePoll/internal/slashcmd"
	"github.com/Oni-Men/SchedulePoll/pkg/dateparser"
	"github.com/Oni-Men/SchedulePoll/pkg/emoji"
	"github.com/bwmarrin/discordgo"
)

type PollService struct {
	bot         *bot.Bot
	pollManager *poll.PollManager
}

func NewPollService(b *bot.Bot) *PollService {
	return &PollService{
		bot: b,
	}
}

func (ps *PollService) Init() {
	ps.pollManager = poll.CreatePollManager()
	ps.Bind()
}

func (p *PollService) Bind() {
	p.bot.AddHandler(p.handleInteractionCreate)
	p.bot.AddHandler(p.handleReactionAdd)
	p.bot.AddHandler(p.handleReactionRemove)
}

func (ps *PollService) handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionModalSubmit:
		ps.handleModalSubmit(s, i)
	case discordgo.InteractionMessageComponent:
		ps.handlePollItems(s, i)
	}
}

func (ps *PollService) handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	if data.CustomID != slashcmd.PollCreateModal {
		return
	}

	res := getModalResponse(&data)
	p, err := createPollFromModalResponse(res)
	if err != nil {
		ps.bot.RespondText(i.Interaction, err.Error())
		fmt.Println(err)
		return
	} else if p == nil {
		ps.bot.RespondText(i.Interaction, "cannot create a new poll")
		return
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponsePong,
		})
	}

	ps.pollManager.AddPoll(p)

	embed := poll.PrintPoll(p)
	msg, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embed: embed,
	})
	if err != nil {
		fmt.Println(err)
	}

	//Add reactions to vote
	for k := 0; k < len(p.Columns); k++ {
		emoji := emoji.ABCDEmoji(k)
		err = s.MessageReactionAdd(i.ChannelID, msg.ID, emoji)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func (ps *PollService) handlePollItems(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.MessageComponentData().CustomID != "poll-item" {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})
}

func (ps *PollService) handleReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := ps.getPollFromMessage(msg)
	if p == nil {
		return
	}

	removeInvalidReactions(p, s, m.ChannelID, m.MessageID, m.Emoji.Name)
	p.AddVote(getIndexFromEmoji(m.Emoji.Name))

	syncVoterCount(p, s, m.ChannelID, m.MessageID)

	embed := poll.PrintPoll(p)
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
	if err != nil {
		fmt.Println(err)
	}
}

func (ps *PollService) handleReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.UserID == s.State.User.ID {
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := ps.getPollFromMessage(msg)
	if p == nil {
		return
	}

	removeInvalidReactions(p, s, m.ChannelID, m.MessageID, m.Emoji.Name)
	p.RemoveVote(getIndexFromEmoji(m.Emoji.Name))

	syncVoterCount(p, s, m.ChannelID, m.MessageID)

	embed := poll.PrintPoll(p)
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
	if err != nil {
		fmt.Println(err)
	}
}

func getModalResponse(data *discordgo.ModalSubmitInteractionData) *modal.ModalResponse {
	options := make([]modal.ModalResponseOption, 0, 10)

	// モーダルの結果から投票作成に必要な情報を抜き出す
	for _, c := range data.Components {
		row, ok := c.(*discordgo.ActionsRow)
		if !ok {
			continue
		}

		if len(row.Components) == 0 {
			continue
		}

		input, ok := row.Components[0].(*discordgo.TextInput)
		if !ok {
			continue
		}

		val := input.Value
		switch input.CustomID {
		case "poll-title":
			options = append(options, modal.WithTitle(val))
		case "poll-description":
			options = append(options, modal.WithDesc(val))
		case "poll-due":
			t, err := time.Parse("2006/01/02 15:04", val)
			if err != nil {
				break
			}
			options = append(options, modal.WithDue(t))
		case "poll-date-list":
			options = append(options, modal.WithDateList(val))
		}
	}

	return modal.NewModalResponse(options...)
}

func createPollFromModalResponse(res *modal.ModalResponse) (*poll.Poll, error) {
	p := poll.CreatePoll()
	p.Title = res.Title()
	p.Description = res.Description()
	p.Due = res.ExpireAt()

	dp := dateparser.NewDateParser(res.DateList())
	for dp.HasNext() {
		if next, err := dp.Next(); err == nil && next != nil {
			p.AddColumn(poll.CreateColumn(
				next.Date,
				next.BeginAt,
				next.EndAt,
			))
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (ps *PollService) getPollFromMessage(msg *discordgo.Message) *poll.Poll {
	if len(msg.Embeds) != 1 {
		return nil
	}

	p, err := ps.getPollFromEmbed(msg.Embeds[0])
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return p
}

func (ps *PollService) getPollFromEmbed(embed *discordgo.MessageEmbed) (*poll.Poll, error) {
	var err error
	if !strings.HasPrefix(embed.Description, "#") {
		return nil, errors.New("invalid format")
	}

	id := strings.TrimPrefix(embed.Description, "#")
	p := ps.pollManager.GetPoll(id)

	if p == nil {
		if p, err = poll.ParsePollEmbed(embed); err != nil {
			return nil, err
		}
		ps.pollManager.AddPoll(p)
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
	return int([]rune(emojiID)[0]) - int(emoji.ABCs[0])
}

func isValidReaction(emojiID string, p *poll.Poll) bool {
	n := getABCIndex(emojiID)
	return n >= 0 && n < len(p.Columns)
}
