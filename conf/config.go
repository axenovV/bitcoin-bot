package conf

import (
	// Системные пакеты
	"io/ioutil"
	"path/filepath"
	// Парсер yaml файлов
	"gopkg.in/yaml.v2"
)

type Config struct {
	// Токен телеграм бота
	Token string `yaml:"token"`
	// Разрешенные айдишники чатов
	AllowedChatIds []int `yaml:"allowed_chat_ids"`
	// Ключевые слова для открывания двери
	ConvertCurrencyCommands []string `yaml:"convert_currency_commands"`
}

func GetDefaultConfig() (*Config, error) {
	var yamlFile []byte
	var err error
	filename, _ := filepath.Abs("./app.yml")
	yamlFile, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := yaml.Unmarshal(yamlFile, &conf); err != nil {
		return nil, err
	}
	return &conf, err
}
