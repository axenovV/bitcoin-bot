package conf

import (
	// Системные пакеты
	"fmt"
	"io/ioutil"
	"path/filepath"
	// Парсер yaml файлов
	"gopkg.in/yaml.v2"
)

type Config struct {
	// Токен телеграм бота
	Token string `yaml:"token"`

	WebhookUrl  string `yaml:"webhook_url"`
	WebhookPort string `yaml:"webhook_port"`

	ServerUrl  string `yaml:"server_url"`
	ServerPort string `yaml:"server_port"`

	ApiUrl string `yaml:"api_url"`

	// Ключевые слова для открывания двери
	ConvertCurrencyCommands []string `yaml:"convert_currency_commands"`
}

func (c *Config) GetFullWebHookUrl() string {
	return fmt.Sprintf("%s:%s", c.WebhookUrl, c.WebhookPort)
}

func (c *Config) GetFullServerUrl() string {
	return fmt.Sprintf("%s:%s", c.ServerUrl, c.ServerPort)
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
