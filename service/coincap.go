package service

import (
	"net/http"
	"github.com/patrickmn/go-cache"
	"time"
	"strings"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func RequestCurrency(currency string) (* http.Response, error) {

	url := "http://coincap.io/page/" + strings.ToUpper(currency)
	value, found := c.Get(url)

	if found {
		return value.(*http.Response), nil
	} else {
		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		if err == nil {
			c.Set(url, resp, cache.DefaultExpiration)
		}
		return resp, err
	}

}