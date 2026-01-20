package main

import (
	"fmt"
	"os"
)

// Version is set at build time via -ldflags
var Version = "dev"

func main() {
	// Load configuration (CLI flags > env vars > config file > defaults)
	cfg := LoadConfig()

	// Handle --version flag
	if cfg.Version {
		fmt.Printf("claude-code-hey-you %s\n", Version)
		return
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	// Handle --test flag
	if cfg.Test {
		client := NewNtfyClient(cfg)
		notification := Notification{
			Title: "claude-code-hey-you test",
			Body:  "If you see this, notifications are working!",
		}
		if err := client.Send(notification); err != nil {
			fmt.Fprintf(os.Stderr, "error sending test notification: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Test notification sent successfully!")
		return
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
