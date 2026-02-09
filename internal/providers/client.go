// Provider client for goclit
package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Additional provider types (extending smart_router.go)
const (
	ProviderGoBro  ProviderType = "gobro"
	ProviderOllama ProviderType = "ollama"
	ProviderClaude ProviderType = "claude" // alias for anthropic
)

// Message for chat
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest for API
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

// ChatResponse from API
type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Client manages provider connections
type Client struct {
	providerType ProviderType
	baseURL      string
	apiKey       string
	model        string
	httpClient   *http.Client
}

// Config for client
type Config struct {
	Provider ProviderType
	APIKey   string
	BaseURL  string
	Model    string
	Timeout  time.Duration
}

// NewClient creates a new provider client
func NewClient(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	c := &Client{
		providerType: cfg.Provider,
		apiKey:       cfg.APIKey,
		model:        cfg.Model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}

	// Set default URLs
	switch cfg.Provider {
	case ProviderOpenRouter:
		c.baseURL = "https://openrouter.ai/api/v1"
		if c.model == "" {
			c.model = "anthropic/claude-sonnet-4"
		}
	case ProviderOpenAI:
		c.baseURL = "https://api.openai.com/v1"
		if c.model == "" {
			c.model = "gpt-4o-mini"
		}
	case ProviderGoogle:
		c.baseURL = "https://generativelanguage.googleapis.com/v1beta"
		if c.model == "" {
			c.model = "gemini-2.0-flash"
		}
	case ProviderGoBro:
		c.baseURL = "http://localhost:8080/v1"
		if c.model == "" {
			c.model = "auto"
		}
	case ProviderOllama:
		c.baseURL = "http://localhost:11434/api"
		if c.model == "" {
			c.model = "llama3"
		}
	case ProviderClaude, ProviderAnthropic:
		c.baseURL = "https://api.anthropic.com/v1"
		if c.model == "" {
			c.model = "claude-sonnet-4-20250514"
		}
	}

	if cfg.BaseURL != "" {
		c.baseURL = cfg.BaseURL
	}

	return c
}

// AutoDetect finds the best available provider
func AutoDetect() (*Client, error) {
	// Priority: OpenRouter > Claude > OpenAI > Gemini > GoBro > Ollama

	// Check OpenRouter
	if key := os.Getenv("OPENROUTER_API_KEY"); key != "" {
		return NewClient(Config{
			Provider: ProviderOpenRouter,
			APIKey:   key,
		}), nil
	}

	// Check Claude
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return NewClient(Config{
			Provider: ProviderClaude,
			APIKey:   key,
		}), nil
	}

	// Check OpenAI
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return NewClient(Config{
			Provider: ProviderOpenAI,
			APIKey:   key,
			Model:    "gpt-4o-mini",
		}), nil
	}

	// Check Gemini
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		return NewClient(Config{
			Provider: ProviderGoogle,
			APIKey:   key,
			Model:    "gemini-2.0-flash",
		}), nil
	}

	// Check GoBro local
	if isReachable("http://localhost:8080/health") {
		return NewClient(Config{
			Provider: ProviderGoBro,
		}), nil
	}

	// Check Ollama
	if isReachable("http://localhost:11434/api/tags") {
		return NewClient(Config{
			Provider: ProviderOllama,
		}), nil
	}

	return nil, fmt.Errorf("no provider available. Set OPENROUTER_API_KEY, ANTHROPIC_API_KEY, OPENAI_API_KEY, GEMINI_API_KEY, or start GoBro/Ollama")
}

func isReachable(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

// Chat sends a message and returns the response
func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	switch c.providerType {
	case ProviderOllama:
		return c.chatOllama(ctx, messages)
	default:
		return c.chatOpenAI(ctx, messages)
	}
}

// chatOpenAI uses OpenAI-compatible API (OpenRouter, Claude, GoBro)
func (c *Client) chatOpenAI(ctx context.Context, messages []Message) (string, error) {
	reqBody := ChatRequest{
		Model:    c.model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// OpenRouter specific headers
	if c.providerType == ProviderOpenRouter {
		req.Header.Set("HTTP-Referer", "https://github.com/biodoia/goclit-ai")
		req.Header.Set("X-Title", "goclit")
	}

	// Claude specific headers
	if c.providerType == ProviderClaude {
		req.Header.Set("x-api-key", c.apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from model")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// chatOllama uses Ollama's API
func (c *Client) chatOllama(ctx context.Context, messages []Message) (string, error) {
	type ollamaReq struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
		Stream   bool      `json:"stream"`
	}

	reqBody := ollamaReq{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := c.baseURL + "/chat"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Ollama error %d: %s", resp.StatusCode, string(body))
	}

	type ollamaResp struct {
		Message Message `json:"message"`
	}

	var chatResp ollamaResp
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	return chatResp.Message.Content, nil
}

// ProviderName returns human-readable provider name
func (c *Client) ProviderName() string {
	switch c.providerType {
	case ProviderOpenRouter:
		return "OpenRouter"
	case ProviderOpenAI:
		return "OpenAI"
	case ProviderGoogle:
		return "Gemini"
	case ProviderGoBro:
		return "GoBro"
	case ProviderOllama:
		return "Ollama"
	case ProviderClaude, ProviderAnthropic:
		return "Claude"
	default:
		return string(c.providerType)
	}
}

// Model returns current model
func (c *Client) Model() string {
	return c.model
}
