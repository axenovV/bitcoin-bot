package bothandler

import (
	"log"

	"github.com/axenovv/bitcoin-bot/conf"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
	"github.com/axenovv/bitcoin-bot/models"
)

const (
	Start = "/start"
	Price = "/price"
	Best = "/best"
	Capitalization = "/cap"
	Help = "/help"
	Portfolio = "/portfolio"

	//coins
	Bitcoin = "bitcoin"
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

	b.BotUpdates.Handle(Start, func(m *tb.Message) {

	})

	b.BotUpdates.Handle(Best, func(m *tb.Message) {
		response, err := models.RequestTopCurrencies()
		if err == nil {
			b.BotUpdates.Send(m.Sender, response.TopCurrencies())
		} else {
			b.BotUpdates.Send(m.Sender, err.Error())
		}
	})

	b.BotUpdates.Handle(Capitalization, func(m *tb.Message) {
		response, err := models.RequestCurrencies(m.Payload)
		if err == nil {
			b.BotUpdates.Send(m.Sender, response.MarketCapCurrencies())
		} else {
			b.BotUpdates.Send(m.Sender, err.Error())
		}
	})

	b.BotUpdates.Handle(Help, func(m *tb.Message) {

	})

	b.BotUpdates.Handle(Price, func(m *tb.Message) {
		var currency string
		if len(m.Payload) == 0 {
			currency = Bitcoin
		} else {
			currency = m.Payload
		}
		response, err := models.RequestCurrencies(currency)
		if err == nil {
			b.BotUpdates.Send(m.Sender, response.GetCurrenciesText())
		} else {
			b.BotUpdates.Send(m.Sender, err.Error())
		}
	})

	b.BotUpdates.Handle(Portfolio, func(m *tb.Message) {

	})

}
