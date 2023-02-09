package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Oni-Men/SchedulePoll/internal/bot"
	"github.com/Oni-Men/SchedulePoll/internal/model/modal"
	"github.com/Oni-Men/SchedulePoll/internal/model/poll"
	"github.com/Oni-Men/SchedulePoll/internal/slashcmd"
	"github.com/Oni-Men/SchedulePoll/pkg/dateparser"
	"github.com/Oni-Men/SchedulePoll/pkg/emoji"
	"github.com/Oni-Men/SchedulePoll/pkg/timeutil"
	"github.com/bwmarrin/discordgo"
	"github.com/go-playground/validator/v10"
)

const (
	INFO_COLOR = 0x66ccff
	ERR_COLOR  = 0xff6666
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
	ps.BindHandlers()
}

func (p *PollService) BindHandlers() {
	p.bot.AddHandler(p.handleInteractionCreate)
	p.bot.AddHandler(p.handleReactionAdd)
	p.bot.AddHandler(p.handleReactionRemove)
}

func (ps *PollService) handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionModalSubmit:
		go ps.handleModalSubmit(s, i)
	}
}

func (ps *PollService) handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "アンケートを作成中...",
		},
	})
	if err != nil {
		log.Println("error responding to modal submi: " + err.Error())
	}

	if data.CustomID != slashcmd.PollCreateModal {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "エラー",
				Description: "不明なモーダルです",
				Color:       ERR_COLOR,
			}},
		})
		return
	}

	res, err := toModalResponse(&data)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "エラー",
				Description: err.Error(),
				Color:       ERR_COLOR,
			}},
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(res); err != nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "エラー",
				Description: err.Error(),
				Color:       ERR_COLOR,
			}},
		})
		return
	}

	if err != nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "エラー",
				Description: err.Error(),
				Color:       ERR_COLOR,
			}},
		})
		return
	}

	p := newPollFromModalResponse(res)
	ps.pollManager.AddPoll(p)
	embed := PrintPoll(p)

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("@here, <@%s> がスケジュール調整アンケートを作成しました", i.Interaction.Member.User.ID),
		Embeds:  []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		go s.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
			Title:       "エラー",
			Description: err.Error(),
			Color:       ERR_COLOR,
		})
		return
	}

	// Add reactions to each vote. Reactions are regional indicators.
	// Reactions count is the same as columns count in the poll.
	// An order of the reactions should be assured. therfore, adding reactions are not going to use goroutine.
	for k := 0; k < len(p.Columns); k++ {
		emoji := emoji.ABCDEmoji(k)
		err := s.MessageReactionAdd(i.ChannelID, msg.ID, emoji)
		if err != nil {
			go s.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
				Title:       "エラー",
				Description: err.Error(),
				Color:       ERR_COLOR,
			})
		}
	}
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

	embed := PrintPoll(p)
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

	embed := PrintPoll(p)
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
	if err != nil {
		fmt.Println(err)
	}
}

// This function converts response data which is discordgo component to our ModalResponse.
// ModalResponse will directly be used for creating a new poll.
func toModalResponse(data *discordgo.ModalSubmitInteractionData) (*modal.ModalResponse, error) {
	options := make([]modal.ModalResponseOption, 0, 10)

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
			if val == "" {
				options = append(options, modal.WithDue(timeutil.GetZeroTime()))
			} else if due, err := parseDue(val); due != nil {
				options = append(options, modal.WithDue(*due))
			} else if err != nil {
				return nil, err
			}
		case "poll-date-list":
			if dateList, err := toDateList(val); err == nil {
				options = append(options, modal.WithDateList(dateList))
			} else {
				return nil, err
			}
		}
	}

	return modal.NewModalResponse(options...), nil
}

// Parse date list string with DateParser.
func toDateList(input string) ([]dateparser.ParsedDateResult, error) {
	list := make([]dateparser.ParsedDateResult, 0, 10)
	dp := dateparser.NewDateParser(input)
	for dp.HasNext() {
		if next, err := dp.Next(); err == nil {
			list = append(list, *next)
		} else {
			return nil, err
		}
	}
	return list, nil
}

func parseDue(input string) (*time.Time, error) {
	res, err := dateparser.ParseInlineDate(input)
	if err == nil {
		due := res.Date.Add(res.BeginAt)
		return &due, nil
	}
	return nil, err
}

func newPollFromModalResponse(res *modal.ModalResponse) *poll.Poll {
	p := poll.CreatePoll()
	p.Title = res.Title()
	p.Description = res.Description()
	p.Due = res.Due()

	for _, date := range res.DateList() {
		p.AddColumn(poll.CreateColumn(
			date.Date,
			date.BeginAt,
			date.EndAt,
		))
	}
	return p
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
		if p, err = ParsePollEmbed(embed); err != nil {
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
