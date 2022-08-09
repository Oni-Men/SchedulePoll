package bot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

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

func Start(s *discordgo.Session, gid string) int {
	s.AddHandler(HandleMessageCreate)
	s.AddHandler(HandleReactionAdd)
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
	parsed, err := ParseSchedule(m.Content)
	if err != nil {
		return
	}

	s.ChannelMessageDelete(m.ChannelID, m.ID)

	formatted := FormatShecule(parsed)

	var c int
	b := NewEmbedBuilder()
	b.Title(":calendar_spiral: 予定投票 :calendar_spiral:")
	for year, list := range *formatted {
		var value string
		for _, t := range list {
			emoji := RegionalIndicators[c]
			value += emoji + " **" + t.Format("01/02") + "**\n\n"
			c++
		}

		b.AddField(&discordgo.MessageEmbedField{
			Name:   strconv.Itoa(year) + "年",
			Value:  value,
			Inline: false,
		})
	}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, b.Build(s))
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < c; i++ {
		emoji := RegionalIndicators[i]
		err = s.MessageReactionAdd(m.ChannelID, msg.ID, emoji)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ParseSchedule(content string) ([]time.Time, error) {
	parts := strings.Split(content, ":")
	if len(parts) != 2 {
		return []time.Time{}, errors.New("invalid format #1")
	}

	if parts[0] != "!予定投票" {
		return []time.Time{}, errors.New("invalid format #2")
	}

	inputs := strings.Split(parts[1], ",")
	schedules := make([]time.Time, 0, len(inputs))

	var t time.Time
	var err error = nil
	year, month, day := time.Now().Date()
	for _, input := range inputs {
		split := strings.Split(strings.TrimSpace(input), "/")

		if len(split) == 3 {
			if year, err = strconv.Atoi(split[0]); err != nil {
				break
			}
			split = split[1:]
		}

		if len(split) == 2 {
			_month, err := strconv.Atoi(split[0])
			if err != nil {
				break
			}
			month = time.Month(_month)
			split = split[1:]
		}

		if len(split) == 1 {
			if day, err = strconv.Atoi(split[0]); err != nil {
				break
			}
		}

		t = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		schedules = append(schedules, t)
	}

	return schedules, err
}

func FormatShecule(list []time.Time) *map[int][]time.Time {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Before(list[j])
	})

	res := make(map[int][]time.Time)
	for _, t := range list {
		y := t.Year()
		dateList := res[y]
		if dateList == nil {
			dateList = make([]time.Time, 0, 10)
		}
		res[y] = append(dateList, t)
	}
	return &res
}

func HandleReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

}
