package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kormiltsev/tbot-chatgpt/pkg/chatgpt"
)

func main() {
	apiKey := os.Getenv("CHATGPT_API_TOKEN")
	if apiKey == "" {
		panic("CHATGPT_API_TOKEN is required")
	}
	model := "gpt-4o-mini"
	client := chatgpt.NewClient(apiKey, model)

	messages := []chatgpt.Message{
		{Role: "system", Content: "Response as a psychotherapist."},
		{Role: "user", Content: "Tell me I'm OK."},
	}

	response, err := client.Chat(messages, 0.7, 100)
	if err != nil {
		log.Fatalf("Chat error: %v", err)
	}

	fmt.Println("ChatGPT:", response)
}
