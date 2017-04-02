package main

import (
	"log"

	"github.com/axenovv/bitcoin-bot/bothandler"
)

func main() {
	startListeningBotUpdate()
}

func startListeningBotUpdate() {
	bot := &bothandler.BotTelegramHandler{}
	err := bot.ConnectToBot()
	if err != nil {
		log.Print(err)
	}
}
