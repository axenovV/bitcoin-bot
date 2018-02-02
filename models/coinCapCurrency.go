package models

import (
	"fmt"
	"github.com/leekchan/accounting"
)

var ac = accounting.Accounting{Symbol: "$", Precision: 2}

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

func (c *Currency) GetMarketCap() string {
	return fmt.Sprintf("%s Rank: %s\n" + "Marketcap: %s", c.Id, c.Rank, c.getMarketCapPrice())
}

func (c *Currency) Print() string {
	return fmt.Sprintf("💵 *%s %s* \n" +
		"📈 Capitalization 24h: %.2f%% \n" +
		"📈 Vol: %s",
		c.Id,
		c.getUsdPrice(),
		c.CapChange24h,
		c.getVolume24())
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

// Private

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
