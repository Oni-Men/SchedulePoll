package main

import (
	"os"

	"github.com/Oni-Men/SchedulePoll/internal/bot"
	"github.com/Oni-Men/SchedulePoll/internal/inits"
)

func goMain() int {
	b := bot.Create(os.Getenv("DISCORD_BOT_TOKEN"))
	if b == nil {
		return 1
	}

	inits.Inititalize(b)

	return b.Start()
}

func main() {
	status := goMain()
	os.Exit(status)
}
