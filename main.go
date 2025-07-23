package main

// TODO: document
// TODO: git init and push
// TODO: deploy on homer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"log"

	godotenv "github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Check interval
	interval = 300 * time.Second
)

// sends an alert message to a Telegram chat
func sendTelegramNotification(bot *tgbotapi.BotAPI, message string) {
	// Telegram Chat ID
	var chatID, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)

	msg := tgbotapi.NewMessage(chatID, message)

	_, err := bot.Send(msg)

	if err != nil {
		fmt.Println("Failed to send Telegram message:", err)
	}
}

// performs an HTTP request and sends an alert if a 4xx error is detected
func checkWebsite(bot *tgbotapi.BotAPI) {
	// URL to monitor
	var urlToCheck = os.Getenv("URL_TO_CHECK")
	
	resp, err := http.Get(urlToCheck)
	if err != nil {
		fmt.Println("Failed to perform HTTP request:", err)
		sendTelegramNotification(bot, fmt.Sprintf("⚠️ Request failed: %v", err))
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to read body from request:", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		message := fmt.Sprintf("🚨 HTTP %d Error Detected!\nURL: %s", resp.StatusCode, urlToCheck)
		sendTelegramNotification(bot, message)
	} else {
		fmt.Printf("✅ %s is OK (%d)\n", urlToCheck, resp.StatusCode)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Telegram Bot Token
	var botToken = os.Getenv("BOT_TOKEN")
	// Initialize Telegram Bot
	fmt.Println("Connecting to bot: ", botToken)
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Failed to initialize Telegram bot:", err)
		os.Exit(1)
	}

	fmt.Println("🔍 Monitoring started...")

	// Testing Telegram gets message
	sendTelegramNotification(bot, "🔍 Monitoring started...")

	// Run the monitoring loop
	for {
		checkWebsite(bot)
		time.Sleep(interval)
	}
}
