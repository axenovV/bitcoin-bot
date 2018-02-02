package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/patrickmn/go-cache"
	"github.com/leekchan/accounting"
	"time"
	"strings"
)

const (
	CoinCapPageUrl = "http://coincap.io/page/"
	CoinMarketCapGlobalUrl = "https://api.coinmarketcap.com/v1/global/"
	CoinMarketcapTopCurrency = "https://api.coinmarketcap.com/v1/ticker/?limit=10"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func RequestCurrencies(currency string) (*Currency, error) {

	url := CoinCapPageUrl + strings.ToUpper(currency)

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
