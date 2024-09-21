package main

import (
	"sseBot/bot"
	"sseBot/config"
	"sseBot/sseapi"
	"sseBot/variable"

	"fmt"
	"log"
	"os"
	"time"

	"os/signal"
	"syscall"
)

func initLog() {
	logFile, err := os.OpenFile("ssebot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open error log file:", err)
		return
	}
	log.SetOutput(logFile)

	// 创建一个信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 来监听信号
	go func() {
		sig := <-sigChan
		log.Printf("接收到信号: %v", sig)
		// 在这里可以添加其他清理操作
		os.Exit(0)
	}()
}

func main() {
	config.InitConfig()
	config := config.GetConfig()
	initLog()
	intChannel := make(chan int)
	postChannel := make(chan variable.Post)
	go bot.InitBot(&config, intChannel, postChannel)

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
