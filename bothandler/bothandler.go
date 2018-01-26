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
	ICO = "/ico"
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
		b.BotUpdates.Send(m.Sender, models.Localization.StartCommand, tb.ModeMarkdown)
	})

	b.BotUpdates.Handle(Best, func(m *tb.Message) {
		response, err := models.RequestTopCurrencies()
		if err == nil {
			b.BotUpdates.Send(m.Sender, response.TopCurrencies(), tb.ModeMarkdown)
		} else {
			b.BotUpdates.Send(m.Sender, err.Error())
		}
	})

	b.BotUpdates.Handle(Capitalization, func(m *tb.Message) {
		if len(m.Payload) == 0 {
			response, err := models.RequestGlobalData()
			if err == nil {
				b.BotUpdates.Send(m.Sender, response.RenderChatMessage())
			} else {
				b.BotUpdates.Send(m.Sender, err.Error())
			}
		} else {
			response, err := models.RequestCurrencies(m.Payload)
			if err == nil {
				b.BotUpdates.Send(m.Sender, response.MarketCapCurrencies())
			} else {
				b.BotUpdates.Send(m.Sender, err.Error())
			}
		}
	})

	b.BotUpdates.Handle(ICO, func(m *tb.Message) {
		response, err := models.GetLiveIcos()
		log.Print(response.Ico.RenderChatMessage())
		if err == nil {
			b.BotUpdates.Send(m.Sender, response.Ico.RenderChatMessage(), tb.ModeMarkdown)
		}
	})

	b.BotUpdates.Handle(Help, func(m *tb.Message) {
		b.BotUpdates.Send(m.Sender, models.Localization.HelpCommand, tb.ModeMarkdown)
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
			b.BotUpdates.Send(m.Sender, response.GetCurrenciesText(), tb.ModeMarkdown)
		} else {
			b.BotUpdates.Send(m.Sender, err.Error())
		}
	})

	b.BotUpdates.Handle(Portfolio, func(m *tb.Message) {

	})

}
