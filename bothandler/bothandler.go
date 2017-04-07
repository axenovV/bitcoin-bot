package bothandler

import (
	"log"
	"net/http"

	"github.com/axenovv/bitcoin-bot/conf"
	"github.com/axenovv/bitcoin-bot/models"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotTelegramHandler struct {
	BotUpdates *tgbotapi.BotAPI
}

func (b *BotTelegramHandler) ConnectToBot(config *conf.Config) error {
	err := b.getBot(config.Token)
	if err != nil {
		log.Print(err)
	}

	b.BotUpdates.Debug = true

	log.Printf("Authorized on account %s", b.BotUpdates.Self.UserName)

	_, err = b.BotUpdates.SetWebhook(tgbotapi.NewWebhook(config.GetFullWebHookUrl() + "/" + b.BotUpdates.Token))
	if err != nil {
		log.Fatal(err)
	}

	updates := b.BotUpdates.ListenForWebhook("/" + b.BotUpdates.Token)

	go http.ListenAndServe(config.GetFullServerUrl(), nil)

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
	case "btc":
		response, err := models.RequestCurrencies("bitcoin")
		if err == nil {
			message.Text = response.GetCurrenciesText()
		}
	case "eth":
		response, err := models.RequestCurrencies("ethereum")
		if err == nil {
			message.Text = response.GetCurrenciesText()
		}
	case "rpl":
		response, err := models.RequestCurrencies("ripple")
		if err == nil {
			message.Text = response.GetCurrenciesText()
		}
	case "dash":
		response, err := models.RequestCurrencies("dash")
		if err == nil {
			message.Text = response.GetCurrenciesText()
		}
	case "litecoin":
		response, err := models.RequestCurrencies("litecoin")
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
