package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// NtfyClient handles sending notifications to ntfy
type NtfyClient struct {
	Server   string
	Topic    string
	Priority string
	client   *http.Client
}

// NewNtfyClient creates a new ntfy client from config
func NewNtfyClient(cfg Config) *NtfyClient {
	return &NtfyClient{
		Server:   strings.TrimSuffix(cfg.Server, "/"),
		Topic:    cfg.Topic,
		Priority: cfg.Priority,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send sends a notification to ntfy
func (c *NtfyClient) Send(n Notification) error {
	url := fmt.Sprintf("%s/%s", c.Server, c.Topic)

	req, err := http.NewRequest("POST", url, strings.NewReader(n.Body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Title", n.Title)
	req.Header.Set("Priority", c.Priority)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy returned status %d", resp.StatusCode)
	}

	return nil
}
