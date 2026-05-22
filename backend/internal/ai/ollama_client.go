package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
	Model   string
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func (c *Client) Generate(prompt string) (string, error) {
	payload, _ := json.Marshal(GenerateRequest{Model: c.Model, Prompt: prompt, Stream: false})
	resp, err := http.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out := GenerateResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Response, nil
}
