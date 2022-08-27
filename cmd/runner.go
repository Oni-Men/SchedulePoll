package main

import (
	"os"

	"github.com/Oni-Men/SchedulePoll/bot"
)

func goMain() int {
	s := bot.Create(os.Getenv("DISCORD_BOT_TOKEN"))
	if s == nil {
		return 1
	}

	return bot.Start(s)
}

func main() {
	status := goMain()
	os.Exit(status)
}
