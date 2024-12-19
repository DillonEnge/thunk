package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	domain     string
	httpClient http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type chatCompletionPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type CompletionResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

type ChatCompletionResponse struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

func NewClient(domain string, httpClient *http.Client) *Client {
	if httpClient == nil {
		return &Client{
			domain: domain,
		}
	}

	return &Client{
		domain:     domain,
		httpClient: *httpClient,
	}
}

func (c *Client) Completion(model string, prompt string) (*http.Response, error) {
	cp := completionPayload{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cp); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/api/generate", c.domain), "application/json", &b)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ChatCompletion(model string, messages []Message) (*http.Response, error) {
	cp := chatCompletionPayload{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cp); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/api/chat", c.domain), "application/json", &b)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
