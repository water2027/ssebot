package config

import (
	"sync"
	"encoding/json"
	"io"
	"log"
	"os"
)

type BotConfig struct {
	TargetGroupName1 string `json:"targetGroupName1"`
	TargetGroupName2 string `json:"targetGroupName2"`
	TimeInterval    int    `json:"timeInterval"`
	Telephone       string `json:"telephone"`
	Email           string `json:"email"`
	Password        string `json:"password"` //用go没复现出来，只能手动复制了
	Str             string `json:"str"`
}

var config BotConfig
var mu sync.Mutex

func InitConfig() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
}

func GetConfig() BotConfig {
	mu.Lock()
	defer mu.Unlock()
	return config
}