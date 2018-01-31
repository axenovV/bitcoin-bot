package main

import (
	"log"

	"github.com/axenovv/bitcoin-bot/bothandler"
	"github.com/axenovv/bitcoin-bot/conf"
)

var handler = bothandler.BotTelegramHandler{}

func main() {
	startListeningBotUpdate()
}

func startListeningBotUpdate() {
	config, configErr := conf.GetDefaultConfig()

	if configErr != nil {
		log.Panic(configErr)
	}


	err := handler.ConnectToBot(config)
	if err != nil {
		log.Panic(err)
	}
}
