package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model   string `json:"model"`
	Content string `json:"content"`
}

type Provider interface {
	Chat(messages []Message) (*Response, error)
}

type OpenAIProvider struct {
	BaseURL string
	APIKey  string
	Model   string
	client  *http.Client
}

func NewOpenAIProvider(baseURL string, apiKey string, model string) *OpenAIProvider {
	return &OpenAIProvider{
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
		Model:   model,
		client:  &http.Client{Timeout: 20 * time.Second},
	}
}

func (p *OpenAIProvider) Chat(messages []Message) (*Response, error) {
	if p.BaseURL == "" || p.APIKey == "" || p.Model == "" {
		return nil, errors.New("base_url, api_key and model are required")
	}

	body := map[string]any{
		"model":    p.Model,
		"messages": messages,
		"stream":   false,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, p.BaseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("provider returned status %d", resp.StatusCode)
	}

	var result struct {
		Model   string `json:"model"`
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	content := ""
	if len(result.Choices) > 0 {
		content = result.Choices[0].Message.Content
	}

	return &Response{Model: result.Model, Content: content}, nil
}
