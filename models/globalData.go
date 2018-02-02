package models

import (
	"fmt"
	"encoding/json"
	)

type ResponseGlobalData struct {
	TotalMarketCapUsd                float64 `json:"total_market_cap_usd"`
	Total24VolumeUsd                 float64 `json:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap     float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies                 int 	 `json:"active_currencies"`
	ActiveAssets               		 int     `json:"active_assets"`
	ActiveMarkets                	 int     `json:"active_markets"`
}

func (r *ResponseGlobalData) GetTotalMarketCapUsd() string {
	return ac.FormatMoney(r.TotalMarketCapUsd)
}

func (r *ResponseGlobalData) GetTotal24VolumeUsd() string {
	return ac.FormatMoney(r.Total24VolumeUsd)
}

func (r *ResponseGlobalData) GetBitcoinPercentageOfMarketCap() string {
	return fmt.Sprintf("%.2f%%", r.BitcoinPercentageOfMarketCap)
}

func (r *ResponseGlobalData) RenderChatMessage() string {
	return fmt.Sprintf("Market Cap: %s \n" +
		"Bitcoin Dominance: %s \n" +
		"24H Volume: %s", r.GetTotalMarketCapUsd(), r.GetBitcoinPercentageOfMarketCap(), r.GetTotal24VolumeUsd())
}

func (r *ResponseGlobalData) UnmurshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &r)
}

// Response Currencies

type ResponseCurrencies struct {
	Result []Currency
}

func (c *ResponseCurrencies) UnmurshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.Result)
}

func (c *ResponseCurrencies) GetCurrenciesText() string {
	text := ""
	for _, value := range c.Result {
		text += value.CurrencyFormating() + "\n"
	}

	return text
}

func (c *ResponseCurrencies) MarketCapCurrencies() string {
	text := ""
	for _, value := range c.Result {
		text += value.GetMarketCap() + "\n"
	}

	return text
}

func (c *ResponseCurrencies) TopCurrencies() string {
	text := ""
	for _, value := range c.Result {
		text += value.SimplePrice() + "\n"
	}

	return text
}