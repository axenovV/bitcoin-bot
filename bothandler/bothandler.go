package bothandler

import (
	"log"

	"github.com/axenovv/bitcoin-bot/conf"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
	"github.com/axenovv/bitcoin-bot/models"
)

const (
	Ticker = "/ticker"
	Portfolio = "/portfolio"
)

type BotTelegramHandler struct {
	BotUpdates *tb.Bot
}

func (b *BotTelegramHandler) ConnectToBot(config *conf.Config) error {

	err := b.getBot(config.Token)

	if err != nil {
		log.Panic(err)
	}

	b.setupHandlers()

	b.BotUpdates.Start()

	return nil
}

func (b *BotTelegramHandler) getBot(token string) error {

	bot, err := tb.NewBot(tb.Settings{
		Token:    token,
		Poller:   &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return err
	}
	b.BotUpdates = bot
	return nil
}

func (b *BotTelegramHandler) setupHandlers() {
	b.BotUpdates.Handle(Ticker, func(m *tb.Message) {
		response, err := models.RequestCurrencies(m.Payload)
		if err == nil {
			b.BotUpdates.Send(m.Sender, response.GetCurrenciesText())
		} else {
			b.BotUpdates.Send(m.Sender, err.Error())
		}
	})
}
