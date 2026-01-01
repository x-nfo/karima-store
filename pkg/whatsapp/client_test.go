package whatsapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Send(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request URL
		if req.URL.String() != "/send" {
			t.Errorf("expected URL /send, got %s", req.URL.String())
		}

		// Test headers
		if req.Header.Get("Authorization") != "test-token" {
			t.Errorf("expected Authorization header test-token, got %s", req.Header.Get("Authorization"))
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
		}

		// Test body
		var body SendMessageRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if body.Target != "08123456789" {
			t.Errorf("expected target 08123456789, got %s", body.Target)
		}
		if body.Message != "Hello World" {
			t.Errorf("expected message Hello World, got %s", body.Message)
		}

		// Send response
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"status": true, "msg": "sent"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL // Override base URL for testing

	// Test Send method
	err := client.Send("08123456789", "Hello World")
	if err != nil {
		t.Errorf("Send() error = %v, wantErr %v", err, nil)
	}
}
