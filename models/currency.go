package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"github.com/leekchan/accounting"
)

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

var ac = accounting.Accounting{Symbol: "$", Precision: 2}


type Currency struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	UsdPrice             string `json:"price_usd"`
	EurPrice             string `json:"price_eur"`
	Rank 				 string `json:"rank"`
	MarketCapUsd 		 string `json:"market_cap_usd"`
	LastUpdated          string `json:"last_updated"`
	PercentChangeOneHour string `json:"percent_change_1h"`
	PercentChangeOneDay  string `json:"percent_change_24h"`
	PercentChangeOneWeek string `json:"percent_change_7d"`
	Volume 				 string `json:"24h_volume_usd"`
}

func (c *Currency) getMarketCapPrice() string {
	marketCapUsd, _ :=  floatFromString(c.MarketCapUsd)
	return ac.FormatMoney(marketCapUsd)
}

func (c *Currency) getVolume24() string {
	volume24, _ :=  floatFromString(c.Volume)
	return ac.FormatMoney(volume24)
}

func (c *Currency) getUsdPrice() string {
	usdPrice, _ :=  floatFromString(c.UsdPrice)
	return ac.FormatMoney(usdPrice)
}

func (c *Currency) SimplePrice() string {
	changeOneDay, _ :=  floatFromString(c.PercentChangeOneDay)
	return fmt.Sprintf("%s  %s  (%.2f%%)", c.Symbol, c.getUsdPrice(), changeOneDay)
}

func (c *Currency) GetMarketCap() string {
	return fmt.Sprintf("%s Rank: %s\n" + "Marketcap: %s", c.Symbol, c.Rank, c.getMarketCapPrice())
}

func (c *Currency) CurrencyFormating() string {
	changeOneHour, _ := floatFromString(c.PercentChangeOneHour)
	changeOneDay, _ :=  floatFromString(c.PercentChangeOneDay)
	changeOneWeek, _ := floatFromString(c.PercentChangeOneWeek)
	return fmt.Sprintf("💵 *%s %s* \n" +
		                      	"📈 Change 1h: %.2f%% \n" +
						      "📈 Change 1d: %.2f%% \n" +
					          "📈 Change 1w: %.2f%% \n" +
						      "📈 Vol: %s",
						      	c.Symbol,
						      	c.getUsdPrice(),
						      				changeOneHour,
						      				changeOneDay,
						      				changeOneWeek,
						      				c.getVolume24())
}

func floatFromString(value string) (float64, error) {
	percent, err := strconv.ParseFloat(value, 64)
	return percent, err
}

func RequestCurrencies(currency string) (*ResponseCurrencies, error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/" + currency)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var response = &ResponseCurrencies{}
	return response, response.UnmurshalJSON(body)
}

func RequestTopCurrencies() (*ResponseCurrencies, error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/?limit=10")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var response = &ResponseCurrencies{}
	return response, response.UnmurshalJSON(body)
}