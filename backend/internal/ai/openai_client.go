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

type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	HTTP    *http.Client
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type,omitempty"`
	} `json:"error,omitempty"`
}

func (c *Client) Generate(prompt string) (string, error) {
	if strings.TrimSpace(c.APIKey) == "" {
		return "", errors.New("OPENAI_API_KEY is required")
	}

	payload, err := json.Marshal(chatCompletionRequest{
		Model: c.requestModel(),
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: "You generate production-ready DevOps artifacts. Return concise, directly usable file contents.",
			},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.2,
	})
	if err != nil {
		return "", err
	}

	endpoint := strings.TrimRight(c.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := c.HTTP
	if client == nil {
		client = &http.Client{Timeout: 2 * time.Minute}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	out := chatCompletionResponse{}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if out.Error != nil && out.Error.Message != "" {
			return "", fmt.Errorf("openai request failed: %s", out.Error.Message)
		}
		return "", fmt.Errorf("openai request failed with status %d", resp.StatusCode)
	}
	if len(out.Choices) == 0 {
		return "", errors.New("openai response contained no choices")
	}

	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

func (c *Client) requestModel() string {
	model := strings.TrimSpace(c.Model)
	if strings.Contains(c.BaseURL, "api.openai.com") {
		return strings.TrimPrefix(model, "openai/")
	}
	return model
}
