package fonnte

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the Fonnte API client
type Client struct {
	apiToken   string
	apiURL     string
	httpClient *http.Client
}

// NewClient creates a new Fonnte client
func NewClient(apiToken, apiURL string) *Client {
	return &Client{
		apiToken: apiToken,
		apiURL:   apiURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessageRequest represents a request to send a WhatsApp message
type SendMessageRequest struct {
	Target  string `json:"target"`         // Phone number (628xxx format)
	Message string `json:"message"`        // Message content
	Type    string `json:"type,omitempty"` // Message type (text, image, video, etc)
	URL     string `json:"url,omitempty"`  // URL for media messages
}

// SendMessageResponse represents the API response
type SendMessageResponse struct {
	Status  bool        `json:"status"`
	Detail  string      `json:"detail,omitempty"`
	Message string      `json:"message,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

// SendMessage sends a WhatsApp message via Fonnte API
func (c *Client) SendMessage(target, message string) (*SendMessageResponse, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("target", target)
	data.Set("message", message)

	// Create request
	req, err := http.NewRequest("POST", c.apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	return &response, nil
}

// SendMessageWithMedia sends a message with media attachment
func (c *Client) SendMessageWithMedia(target, message, mediaURL, mediaType string) (*SendMessageResponse, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("target", target)
	data.Set("message", message)
	data.Set("url", mediaURL)
	data.Set("type", mediaType)

	// Create request
	req, err := http.NewRequest("POST", c.apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// SendBulkMessage sends messages to multiple recipients
func (c *Client) SendBulkMessage(targets []string, message string) ([]SendMessageResponse, error) {
	// Join targets with comma
	targetStr := strings.Join(targets, ",")

	// Prepare JSON body for bulk
	jsonBody, err := json.Marshal(map[string]string{
		"target":  targetStr,
		"message": message,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response - bulk returns array
	var responses []SendMessageResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		// Try single response
		var response SendMessageResponse
		if err2 := json.Unmarshal(body, &response); err2 == nil {
			return []SendMessageResponse{response}, nil
		}
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return responses, nil
}

// GetDeviceStatus checks the WhatsApp device status
func (c *Client) GetDeviceStatus() (*SendMessageResponse, error) {
	statusURL := strings.Replace(c.apiURL, "/send", "/device", 1)

	req, err := http.NewRequest("POST", statusURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
