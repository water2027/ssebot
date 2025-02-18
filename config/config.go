package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Webhook      string `yaml:"webhook" json:"webhook"`
	TemplateStr  string `yaml:"template_str" json:"template_str"`
	TimeInterval int64  `yaml:"time_interval" json:"time_interval"`
	Telephone    string `yaml:"telephone" json:"telephone"`
	Email        string `yaml:"email" json:"email"`
	Password     string `yaml:"password" json:"password"`
	StartNum     int    `yaml:"startNum" json:"startNum"`
}

var BotConfig Config
var NowNum = 0

func InitConfig() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	err = yaml.Unmarshal(byteValue, &BotConfig)
	if err != nil {
		log.Println(err)
		return
	}
	NowNum = BotConfig.StartNum
}
