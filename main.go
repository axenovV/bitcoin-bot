package main

import (
	"log"

	"github.com/axenovv/bitcoin-bot/bothandler"
	"github.com/axenovv/bitcoin-bot/conf"
)

func main() {
	startListeningBotUpdate()
}

func startListeningBotUpdate() {
	config, configErr := conf.GetDefaultConfig()
	if configErr != nil {
		log.Print(configErr)
	}
	bot := &bothandler.BotTelegramHandler{}
	err := bot.ConnectToBot(config)
	if err != nil {
		log.Print(err)
	}
}
