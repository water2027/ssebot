package main

import (
	"sseBot/config"
	"sseBot/sseapi"
	"sseBot/bot"

	"fmt"
	"log"
	"os"
	"time"
)

func initLog() {
	logFile, err := os.OpenFile("ssebot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open error log file:", err)
		return
	}
	log.SetOutput(logFile)
}

func main() {
	config.InitConfig()
	config := config.GetConfig()
	initLog()
	intChannel := make(chan int)
	postChannel := make(chan sseapi.Post)	
	go bot.InitBot(&config, intChannel,postChannel)

	ticker := time.NewTicker(time.Duration(config.TimeInterval) * time.Minute)
	for {
		select {
		case <-intChannel:
			return
		case <-ticker.C:
			sseapi.GetPosts(postChannel, &config)
		}
	}

}
