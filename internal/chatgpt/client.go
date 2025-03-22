package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ChatGPTClient struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func NewClient(apiKey, model string) *ChatGPTClient {
	return &ChatGPTClient{
		APIKey:  apiKey,
		BaseURL: "https://api.openai.com/v1/chat/completions",
		Model:   model,
		Client:  &http.Client{},
	}
}

func (c *ChatGPTClient) Chat(messages []Message, temperature float32, maxTokens int) (string, error) {
	reqBody := chatRequest{
		Model:       c.Model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}
	resp, err := c.sendRequest(&reqBody)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *ChatGPTClient) sendRequest(request *chatRequest) (*chatResponse, error) {
	reqJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var res chatResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Choices) == 0 {
		return nil, errors.New("no response from ChatGPT")
	}

	return &res, nil
}
