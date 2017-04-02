package main

func main() {
	startListeningBotUpdate()
}

func startListeningBotUpdate() {
	bot := &BotTelegramHandler{}
	err := bot.ConnectToBot()
}
