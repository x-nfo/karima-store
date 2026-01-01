package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	token   string
	baseURL string
	client  *http.Client
}

type SendMessageRequest struct {
	Target  string `json:"target"`
	Message string `json:"message"`
}

func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: "https://api.fonnte.com",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Send(to string, message string) error {
	reqBody := SendMessageRequest{
		Target:  to,
		Message: message,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/send", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API returned error status: %d", resp.StatusCode)
	}

	return nil
}
