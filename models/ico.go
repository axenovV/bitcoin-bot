package models

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Ico struct {
	Ico IcoResponse `json:"ico"`
}

type IcoResponse struct {
	Live []CompanyIco      `json:"live"`
	Upcoming []CompanyIco  `json:"upcoming"`
	Finished []CompanyIco  `json:"finished"`
}

func (ico *IcoResponse) GetFirstTenLiveIcos() []CompanyIco {
	if len(ico.Live) <= 10 {
		return ico.Live
	} else {
		 b := make([]CompanyIco, 0)

		for i := 0; i < 10; i++ {
		    b = append(b, ico.Live[i])
		}
		return b
	}
}
func (ico *IcoResponse) RenderChatMessage() string {
	text := ""
	for _, value := range ico.GetFirstTenLiveIcos() {
		text += value.RenderChatMessage() + "\n\n"
	}

	return text
}

type CompanyIco struct  {
	Name string        `json:"name"`
	Description string `json:"description"`
	WebsiteLink string `json:"website_link"`
	StartTime string   `json:"start_time"`
	EndTime string     `json:"end_time"`
}

func (ico *CompanyIco) RenderChatMessage() string {
	return fmt.Sprintf("*%s* \nDescription: %s\nWebsite: %s", ico.Name, ico.Description, ico.WebsiteLink)
}

func GetLiveIcos() (*Ico, error) {
	url := "https://api.icowatchlist.com/public/v1/live"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var response = &Ico{}

	return response, json.Unmarshal(body, response)
}