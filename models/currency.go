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
	"strings"
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
	Id                   string  `json:"id"`
	Name                 string  `json:"display_name"`
	Rank  				 int8    `json:"rank"`
	UsdPrice             float64 `json:"price_usd"`
	EurPrice             float64 `json:"price_eur"`
	MarketCapUsd 		 float64 `json:"market_cap"`
	Volume 				 float64 `json:"volume"`
	CapChange24h         float64 `json:"cap24hrChange"`
	}

func (c *Currency) getMarketCapPrice() string {
	return ac.FormatMoney(c.MarketCapUsd)
}

func (c *Currency) getVolume24() string {
	return ac.FormatMoney(c.Volume)
}

func (c *Currency) getUsdPrice() string {
	return ac.FormatMoney(c.UsdPrice)
}

func (c *Currency) SimplePrice() string {
	return fmt.Sprintf("%s  %s  (%.2f%%)", c.Id, c.getUsdPrice(), c.CapChange24h)
}

func (c *Currency) GetMarketCap() string {
	return fmt.Sprintf("%s Rank: %s\n" + "Marketcap: %s", c.Id, c.Rank, c.getMarketCapPrice())
}

func (c *Currency) CurrencyFormating() string {

	return fmt.Sprintf("💵 *%s %s* \n" +
		                      	"📈 Capitalization 24h: %.2f%% \n" +
						      "📈 Vol: %s",
						      	c.Id,
						      	c.getUsdPrice(),
						      		c.CapChange24h,
						      				c.getVolume24())
}

// Other

func floatFromString(value string) (float64, error) {
	percent, err := strconv.ParseFloat(value, 64)
	return percent, err
}

func RequestCurrencies(currency string) (*Currency, error) {

	url := "http://coincap.io/page/" + strings.ToUpper(currency)

	value, found := c.Get(url)

	if found {
		return value.(*Currency), nil
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

		var response = &Currency{}
		err = json.Unmarshal(body, &response)
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