package main

import (
	"fmt"
	"os"
)

func main() {
	// Load configuration (CLI flags > env vars > config file > defaults)
	cfg := LoadConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	// Parse hook event from stdin
	event, err := ParseHookEvent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing stdin: %s\n", err)
		os.Exit(1)
	}

	// Format the notification based on the event type
	notification := FormatNotification(event)

	// Skip if notification should be filtered (e.g., duplicate permission notifications)
	if notification == nil {
		return
	}

	// Send to ntfy
	client := NewNtfyClient(cfg)
	if err := client.Send(*notification); err != nil {
		fmt.Fprintf(os.Stderr, "error sending notification: %s\n", err)
		os.Exit(1)
	}
}
