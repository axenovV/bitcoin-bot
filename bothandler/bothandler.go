package bothandler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/axenovv/bitcoin-bot/models"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotTelegramHandler struct {
	BotUpdates *tgbotapi.BotAPI
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
func (b *BotTelegramHandler) ConnectToBot() error {
	err := b.getBot("363505789:AAFQz5eq5oEgNYWfe5J0HEal_IGXeyuT8lM")
	if err != nil {
		log.Print(err)
	}

	b.BotUpdates.Debug = true

	log.Printf("Authorized on account %s", b.BotUpdates.Self.UserName)

	_, err = b.BotUpdates.SetWebhook(tgbotapi.NewWebhook("https://webhook.vkprism.ru:80/" + b.BotUpdates.Token))
	if err != nil {
		log.Fatal(err)
	}

	updates := b.BotUpdates.ListenForWebhook("/" + b.BotUpdates.Token)

	go http.ListenAndServe("95.213.251.26:8443", nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			b.handleCommand(update.Message.Command(), msg)
		}
	}
	return nil
}

func (b *BotTelegramHandler) handleCommand(command string, message tgbotapi.MessageConfig) {

	switch command {
	case "usd":
		response, err := models.RequestCurrencies()
		if err == nil {
			message.Text = response.GetCurrenciesText()
		}
	default:
	}
	b.BotUpdates.Send(message)
}

func (b *BotTelegramHandler) getBot(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	b.BotUpdates = bot
	return nil
}

func (b *BotTelegramHandler) startListenUpdates() (tgbotapi.UpdatesChannel, error) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := b.BotUpdates.GetUpdatesChan(updateConfig)
	return updates, err
}
