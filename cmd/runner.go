package main

import (
	"flag"
	"os"

	"github.com/Oni-Men/SchedulePoll/bot"
)

var Token = flag.String("token", "", "Token of your bot")
var GuildID = flag.String("guild", "", "Id of your guild")

func goMain() int {
	flag.Parse()

	s := bot.Create(*Token)
	if s == nil {
		return 1
	}

	return bot.Start(s, *GuildID)
}

func main() {
	status := goMain()
	os.Exit(status)
}
