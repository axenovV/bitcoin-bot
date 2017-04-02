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

type Currency struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	UsdPrice    string `json:"price_usd"`
	EurPrice    string `json:"price_eur"`
	LastUpdated string `json:"last_updated"`
}

func (c Currency) CurrencyFormating() string {
	i, err := strconv.ParseInt(c.LastUpdated, 10, 64)
	timestamp := int64(i)
	if err != nil {
		return "parsing time error"
	}
	t := time.Unix(timestamp, 0)
	return fmt.Sprintf("💵 %s: $%s  %s", c.Symbol, c.UsdPrice, t.Format(layout))
}

func (c Currency) AsUsdValue() string {
	return ""
}

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

func RequestCurrencies() (*ResponseCurrencies, error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/?limit=2")
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
