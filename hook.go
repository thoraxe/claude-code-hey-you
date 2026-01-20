package main

import (
	"encoding/json"
	"io"
	"os"
)

// HookEvent represents the JSON payload from Claude Code hooks
type HookEvent struct {
	HookEventName    string                 `json:"hook_event_name"`
	ToolName         string                 `json:"tool_name,omitempty"`
	ToolInput        map[string]interface{} `json:"tool_input,omitempty"`
	Message          string                 `json:"message,omitempty"`
	NotificationType string                 `json:"notification_type,omitempty"`
	Cwd              string                 `json:"cwd,omitempty"`
}

// Notification represents a formatted notification to send
type Notification struct {
	Title string
	Body  string
}

// ParseHookEvent reads and parses the hook event from stdin
func ParseHookEvent() (*HookEvent, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	// Handle empty input gracefully
	if len(data) == 0 {
		return &HookEvent{}, nil
	}

	var event HookEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// FormatNotification converts a HookEvent into a Notification with appropriate title and body
func FormatNotification(event *HookEvent) Notification {
	switch event.HookEventName {
	case "PreToolUse":
		return formatPreToolUse(event)
	case "Notification":
		return formatNotificationEvent(event)
	case "Stop":
		return formatStop(event)
	case "SessionStart":
		return formatSessionStart(event)
	default:
		// Generic fallback
		return Notification{
			Title: "Claude Code",
			Body:  event.HookEventName,
		}
	}
}

// formatPreToolUse formats a PreToolUse event based on the tool type
func formatPreToolUse(event *HookEvent) Notification {
	title := "Permission: " + event.ToolName

	var body string
	switch event.ToolName {
	case "Bash":
		// Extract command from tool_input
		if cmd, ok := event.ToolInput["command"].(string); ok {
			body = cmd
		} else if desc, ok := event.ToolInput["description"].(string); ok {
			body = desc
		}
	case "Write", "Edit", "Read":
		// Extract file_path from tool_input
		if path, ok := event.ToolInput["file_path"].(string); ok {
			body = path
		}
	default:
		// For other tools, try to get a meaningful summary
		if data, err := json.Marshal(event.ToolInput); err == nil {
			body = string(data)
			// Truncate if too long
			if len(body) > 200 {
				body = body[:197] + "..."
			}
		}
	}

	if body == "" {
		body = event.ToolName + " operation"
	}

	return Notification{
		Title: title,
		Body:  body,
	}
}

// formatNotificationEvent formats a Notification hook event
func formatNotificationEvent(event *HookEvent) Notification {
	body := event.Message
	if body == "" {
		body = "Notification from Claude Code"
	}

	return Notification{
		Title: "Claude Code",
		Body:  body,
	}
}

// formatStop formats a Stop event
func formatStop(event *HookEvent) Notification {
	body := "Task completed"
	if event.Cwd != "" {
		body = event.Cwd
	}

	return Notification{
		Title: "Task Complete",
		Body:  body,
	}
}

// formatSessionStart formats a SessionStart event
func formatSessionStart(event *HookEvent) Notification {
	body := "New session"
	if event.Cwd != "" {
		body = event.Cwd
	}

	return Notification{
		Title: "Session Started",
		Body:  body,
	}
}
