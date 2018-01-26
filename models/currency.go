package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/patrickmn/go-cache"
	"github.com/leekchan/accounting"
	"time"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

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

var ac = accounting.Accounting{Symbol: "$", Precision: 2}

// Currency

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

// Other

func floatFromString(value string) (float64, error) {
	percent, err := strconv.ParseFloat(value, 64)
	return percent, err
}

func RequestCurrencies(currency string) (*ResponseCurrencies, error) {

	url := "https://api.coinmarketcap.com/v1/ticker/" + currency

	value, found := c.Get(url)

	if found {
		return value.(*ResponseCurrencies), nil
	} else {
		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		var response = &ResponseCurrencies{}
		err = response.UnmurshalJSON(body)
		if err == nil {
			c.Set(url, response, cache.DefaultExpiration)
		}
		return response, err
	}
}

func RequestGlobalData() (*ResponseGlobalData, error)  {

	url := "https://api.coinmarketcap.com/v1/global/"

	value, found := c.Get(url)

	if found {
		return value.(*ResponseGlobalData), nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var response = &ResponseGlobalData{}
	err = response.UnmurshalJSON(body)
	if err == nil {
		c.Set(url, response, cache.DefaultExpiration)
	}
	return response, err
}

func RequestTopCurrencies() (*ResponseCurrencies, error) {

	url := "https://api.coinmarketcap.com/v1/ticker/?limit=10"

	value, found := c.Get(url)

	if found {
		return value.(*ResponseCurrencies), nil
	}
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var response = &ResponseCurrencies{}
	err = response.UnmurshalJSON(body)
	if err == nil {
		c.Set(url, response, cache.DefaultExpiration)
	}
	return response, err
}