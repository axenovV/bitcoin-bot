package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const layout = "Jan 2 2006 3:04pm MST"

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

func (c *Currency) SimplePrice() string {
	changeOneDay, _ :=  floatFromString(c.PercentChangeOneDay)
	return fmt.Sprintf("%s  $%s  (%.2f%%)", c.Symbol, c.UsdPrice, changeOneDay)
}

func (c *Currency) GetMarketCap() string {
	marketCapUsd, _ :=  floatFromString(c.MarketCapUsd)
	return fmt.Sprintf("%s Rank: %s\n" + "Marketcap: %.2f$", c.Symbol, c.Rank, marketCapUsd)
}

func (c *Currency) GetLastUpdateTimeAsString() string {
	i, err := strconv.ParseInt(c.LastUpdated, 10, 64)
	timestamp := int64(i)
	if err != nil {
		return "parsing time error"
	}
	t := time.Unix(timestamp, 0)
	return t.Format(layout)
}

func (c *Currency) CurrencyFormating() string {
	changeOneHour, _ := floatFromString(c.PercentChangeOneHour)
	changeOneDay, _ :=  floatFromString(c.PercentChangeOneDay)
	changeOneWeek, _ := floatFromString(c.PercentChangeOneWeek)
	volume24, _ := floatFromString(c.Volume)
	return fmt.Sprintf("💵 %s $%s \n" +
		                      "🕔 Last Update: %s \n" +
		                      	"📈 Change 1h: %.2f%% \n" +
						      "📈 Change 1d: %.2f%% \n" +
					          "📈 Change 1w: %.2f%% \n" +
						      "📈 Vol: %.0f$",
						      	c.Symbol, c.UsdPrice,
								c.GetLastUpdateTimeAsString(),
									changeOneHour, changeOneDay,
										changeOneWeek,
											volume24)
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