package inits

import (
	"github.com/Oni-Men/SchedulePoll/internal/bot"
	"github.com/Oni-Men/SchedulePoll/internal/service"
	"github.com/Oni-Men/SchedulePoll/internal/slashcmd"
)

func Inititalize(b *bot.Bot) {
	b.AddSlashCommand(new(slashcmd.MeetmeCommand))

	ps := service.NewPollService(b)
	ps.Init()

}
