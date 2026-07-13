package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	ProviderOpenAIChat      = "openai_chat"
	ProviderOpenAI          = "openai"
	ProviderOpenAIResponses = "openai_responses"
	ProviderGemini          = "gemini"
	ProviderAnthropic       = "anthropic"
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

type providerConfig struct {
	Provider string
	BaseURL  string
	APIKey   string
	Model    string
}

func NewProvider(provider string, baseURL string, apiKey string, model string) Provider {
	cfg := providerConfig{
		Provider: provider,
		BaseURL:  strings.TrimRight(baseURL, "/"),
		APIKey:   apiKey,
		Model:    model,
	}

	switch provider {
	case ProviderOpenAIResponses:
		return &OpenAIResponsesProvider{cfg: cfg, client: httpClient()}
	case ProviderGemini:
		return &GeminiProvider{cfg: cfg, client: httpClient()}
	case ProviderAnthropic:
		return &AnthropicProvider{cfg: cfg, client: httpClient()}
	default:
		return &OpenAIChatProvider{cfg: cfg, client: httpClient()}
	}
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 20 * time.Second}
}

func (c providerConfig) validate() error {
	if c.BaseURL == "" || c.APIKey == "" || c.Model == "" {
		return errors.New("base_url, api_key and model are required")
	}
	return nil
}

type OpenAIChatProvider struct {
	cfg    providerConfig
	client *http.Client
}

func (p *OpenAIChatProvider) Chat(messages []Message) (*Response, error) {
	if err := p.cfg.validate(); err != nil {
		return nil, err
	}

	body := map[string]any{
		"model":    p.cfg.Model,
		"messages": messages,
		"stream":   false,
	}

	var result struct {
		Model   string `json:"model"`
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err := p.postJSON(p.cfg.BaseURL+"/chat/completions", body, &result); err != nil {
		return nil, err
	}

	content := ""
	if len(result.Choices) > 0 {
		content = result.Choices[0].Message.Content
	}
	return &Response{Model: fallback(result.Model, p.cfg.Model), Content: content}, nil
}

func (p *OpenAIChatProvider) postJSON(url string, body any, out any) error {
	return postJSON(p.client, url, p.cfg.APIKey, body, out)
}

type OpenAIResponsesProvider struct {
	cfg    providerConfig
	client *http.Client
}

func (p *OpenAIResponsesProvider) Chat(messages []Message) (*Response, error) {
	if err := p.cfg.validate(); err != nil {
		return nil, err
	}

	body := map[string]any{
		"model": p.cfg.Model,
		"input": toResponsesInput(messages),
		"store": false,
	}

	var result struct {
		ID         string `json:"id"`
		Model      string `json:"model"`
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := postJSONWithHeaders(p.client, p.cfg.BaseURL+"/responses", p.cfg.APIKey, body, &result, map[string]string{
		"Accept":    "application/json",
		"x-api-key": p.cfg.APIKey,
	}); err != nil {
		return nil, err
	}

	content := result.OutputText
	if content == "" {
		for _, item := range result.Output {
			for _, part := range item.Content {
				if part.Text != "" {
					content += part.Text
				}
			}
		}
	}
	return &Response{Model: fallback(result.Model, p.cfg.Model), Content: content}, nil
}

type GeminiProvider struct {
	cfg    providerConfig
	client *http.Client
}

type AnthropicProvider struct {
	cfg    providerConfig
	client *http.Client
}

func (p *AnthropicProvider) Chat(messages []Message) (*Response, error) {
	if err := p.cfg.validate(); err != nil {
		return nil, err
	}

	body := map[string]any{
		"model":      p.cfg.Model,
		"max_tokens": 256,
		"messages":   toAnthropicMessages(messages),
	}

	var result struct {
		Model   string `json:"model"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := postAnthropicJSON(p.client, p.cfg.BaseURL+"/messages", p.cfg.APIKey, body, &result); err != nil {
		return nil, err
	}

	content := ""
	for _, part := range result.Content {
		if part.Type == "text" || part.Type == "" {
			content += part.Text
		}
	}
	return &Response{Model: fallback(result.Model, p.cfg.Model), Content: content}, nil
}

func (p *GeminiProvider) Chat(messages []Message) (*Response, error) {
	if err := p.cfg.validate(); err != nil {
		return nil, err
	}

	body := map[string]any{
		"contents": toGeminiContents(messages),
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.cfg.BaseURL, p.cfg.Model, p.cfg.APIKey)
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := postJSON(p.client, url, "", body, &result); err != nil {
		return nil, err
	}

	content := ""
	if len(result.Candidates) > 0 {
		for _, part := range result.Candidates[0].Content.Parts {
			content += part.Text
		}
	}
	return &Response{Model: p.cfg.Model, Content: content}, nil
}

func postJSON(client *http.Client, url string, apiKey string, body any, out any) error {
	return postJSONWithHeaders(client, url, apiKey, body, out, nil)
}

func postJSONWithHeaders(client *http.Client, url string, apiKey string, body any, out any, headers map[string]string) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	for key, value := range headers {
		if value != "" {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		detail := strings.TrimSpace(string(body))
		if detail == "" {
			return fmt.Errorf("provider returned status %d", resp.StatusCode)
		}
		return fmt.Errorf("provider returned status %d: %s", resp.StatusCode, detail)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func postAnthropicJSON(client *http.Client, url string, apiKey string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		detail := strings.TrimSpace(string(body))
		if detail == "" {
			return fmt.Errorf("provider returned status %d", resp.StatusCode)
		}
		return fmt.Errorf("provider returned status %d: %s", resp.StatusCode, detail)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func toResponsesInput(messages []Message) []map[string]any {
	input := make([]map[string]any, 0, len(messages))
	for _, message := range messages {
		input = append(input, map[string]any{
			"role": message.Role,
			"content": []map[string]string{
				{"type": "input_text", "text": message.Content},
			},
		})
	}
	return input
}

func toGeminiContents(messages []Message) []map[string]any {
	contents := make([]map[string]any, 0, len(messages))
	for _, message := range messages {
		role := "user"
		if message.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, map[string]any{
			"role": role,
			"parts": []map[string]string{
				{"text": message.Content},
			},
		})
	}
	return contents
}

func toAnthropicMessages(messages []Message) []map[string]string {
	result := make([]map[string]string, 0, len(messages))
	for _, message := range messages {
		role := message.Role
		if role != "assistant" {
			role = "user"
		}
		result = append(result, map[string]string{
			"role":    role,
			"content": message.Content,
		})
	}
	return result
}

func fallback(value string, fallbackValue string) string {
	if value == "" {
		return fallbackValue
	}
	return value
}
