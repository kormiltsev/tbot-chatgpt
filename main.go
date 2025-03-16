package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kormiltsev/tbot-chatgpt/pkg/chatgpt"
)

func main() {
	apiKey := os.Getenv("CHATGPT_API_TOKEN")
	if apiKey == "" {
		panic("CHATGPT_API_TOKEN is required")
	}

	model := "gpt-4o-mini"
	client := chatgpt.NewClient(apiKey, model)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		panic("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message != nil {
			msg := update.Message.Text
			log.Printf("[ %s ] %s", update.Message.From.UserName, update.Message.Text)

			replytext, err := client.NewUserRequest(msg).WithMaxTokens(100).WithTemperature(0.7).Send()
			if err != nil {
				replytext = "error: " + err.Error()
			}
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, replytext)

			if _, err := bot.Send(reply); err != nil {
				log.Println("Failed to send message:", err)
			}
		}
	}
}
