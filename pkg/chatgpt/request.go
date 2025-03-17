package chatgpt

import (
	"log"
)

var (
	defaultTemperature float32 = 0.7
	defaultMaxTokens   int     = 100
)

func (c *ChatGPTClient) NewUserRequest(request string) *ChatRequest {
	newmessage := Message{Role: "user", Content: request}
	reqBody := ChatRequest{
		Model:       c.Model,
		Messages:    []Message{newmessage},
		Temperature: defaultTemperature,
		MaxTokens:   defaultMaxTokens,
		client:      c,
	}
	return &reqBody
}

func (chatreq *ChatRequest) WithTemperature(temperature float32) *ChatRequest {
	chatreq.Temperature = temperature
	if temperature < 0 || temperature > 2 {
		chatreq.Temperature = defaultTemperature
	}
	return chatreq
}

func (chatreq *ChatRequest) WithMaxTokens(maxTokens int) *ChatRequest {
	chatreq.MaxTokens = maxTokens
	if maxTokens <= 0 {
		chatreq.MaxTokens = defaultMaxTokens
	}
	return chatreq
}

func (chatreq *ChatRequest) WithSystemMessage(message string) *ChatRequest {
	if len(message) == 0 {
		return chatreq
	}
	newmessage := make([]Message, len(chatreq.Messages)+1)
	newmessage[0] = Message{Role: "system", Content: message}
	copy(newmessage[1:], chatreq.Messages)
	return chatreq
}

func (chatreq *ChatRequest) Send() (string, error) {
	sysprompts := GetSysPrompts()
	if len(sysprompts) != 0 {
		sysprms := make([]Message, len(sysprompts))
		for i := range sysprompts {
			sysprms[i] = Message{Role: "system", Content: sysprompts[i]}
		}
		chatreq.Messages = append(sysprms, chatreq.Messages...)
	}

	resp, err := chatreq.client.sendRequest(chatreq)
	if err != nil {
		return "", err
	}

	log.Println("[ GPT ] Choises =", len(resp.Choices))
	log.Printf("[ GPT ] (%s):%s", resp.Choices[0].Message.Role, resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}
