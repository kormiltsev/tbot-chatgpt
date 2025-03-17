package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kormiltsev/tbot-chatgpt/configs"
	"github.com/kormiltsev/tbot-chatgpt/internal/tbot"
	"github.com/kormiltsev/tbot-chatgpt/pkg/chatgpt"
)

func main() {

	envtemp := os.Getenv("CHATGPT_API_DEFAULT_TEMPERATURE")
	if envtemp != "" {
		envtempfloat, err := strconv.ParseFloat(envtemp, 32)
		if err != nil {
			log.Println("ENV wrong format", "CHATGPT_API_DEFAULT_TEMPERATURE", envtemp)
		}
		configs.DefaultTemperature = float32(envtempfloat)
	}

	envmaxtok := os.Getenv("CHATGPT_API_MAX_TOKENS")
	if envmaxtok != "" {
		envmaxtokint, err := strconv.Atoi(envmaxtok)
		if err != nil {
			log.Println("ENV wrong format", "CHATGPT_API_MAX_TOKENS", envmaxtok)
		}
		configs.DefaultMaxTokens = envmaxtokint
	}

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
		if update.Message == nil {
			continue
		}

		var replytext string
		if update.Message.IsCommand() {
			replytext, err = tbot.Command(bot, update.Message)
			if err != nil {
				replytext = "internal error" + err.Error()
			}
		} else {
			msg := update.Message.Text
			log.Printf("[ %s ] %s", update.Message.From.UserName, update.Message.Text)

			replytext, err = client.NewUserRequest(msg).WithMaxTokens(configs.DefaultMaxTokens).WithTemperature(configs.DefaultTemperature).Send()
			if err != nil {
				replytext = "error: " + err.Error()
			}
		}
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, replytext)
		if _, err := bot.Send(reply); err != nil {
			log.Println("Failed to send message:", err)
		}

	}
}
