package main

import (
	"os"
	"strings"

	"github.com/Oni-Men/SchedulePoll/internal/bot"
	"github.com/Oni-Men/SchedulePoll/internal/inits"
)

func goMain() int {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	cleanupOption := os.Getenv("CLEANUP_CMDS")
	cleanup := strings.ToUpper(cleanupOption) == "TRUE"

	b := bot.Create(token, cleanup)
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
