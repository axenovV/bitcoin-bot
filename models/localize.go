package models

import (
	"os"
	"fmt"
	"encoding/json"
)

var Localization = LoadLocalization("locale/en-us.all.json")

type Localize struct {

	StartCommand string `json:"start_command"`
	HelpCommand string `json:"help_command"`

}

func LoadLocalization(file string) Localize {
	var config Localize
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}