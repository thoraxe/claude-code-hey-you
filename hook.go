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
// Returns nil if the notification should be skipped (e.g., duplicate permission notifications)
func FormatNotification(event *HookEvent) *Notification {
	switch event.HookEventName {
	case "PreToolUse", "PermissionRequest":
		return formatPermissionRequest(event)
	case "Notification":
		return formatNotificationEvent(event)
	case "Stop":
		return formatStop(event)
	case "SessionStart":
		return formatSessionStart(event)
	default:
		// Generic fallback
		return &Notification{
			Title: "Claude Code",
			Body:  event.HookEventName,
		}
	}
}

// formatPermissionRequest formats a PermissionRequest/PreToolUse event based on the tool type
func formatPermissionRequest(event *HookEvent) *Notification {
	title := "Permission: " + event.ToolName

	var body string
	switch event.ToolName {
	case "Bash":
		// Extract command and description from tool_input
		if cmd, ok := event.ToolInput["command"].(string); ok {
			body = cmd
			// Truncate long commands
			if len(body) > 300 {
				body = body[:297] + "..."
			}
		}
		// Add description if available
		if desc, ok := event.ToolInput["description"].(string); ok && desc != "" {
			if body != "" {
				body = desc + "\n\n" + body
			} else {
				body = desc
			}
		}
	case "Write":
		if path, ok := event.ToolInput["file_path"].(string); ok {
			body = "Create/overwrite: " + path
		}
	case "Edit":
		if path, ok := event.ToolInput["file_path"].(string); ok {
			body = "Edit: " + path
		}
	case "Read":
		if path, ok := event.ToolInput["file_path"].(string); ok {
			body = "Read: " + path
		}
	case "Task":
		// Agent task
		if desc, ok := event.ToolInput["description"].(string); ok {
			body = desc
		}
		if prompt, ok := event.ToolInput["prompt"].(string); ok && body == "" {
			body = prompt
			if len(body) > 200 {
				body = body[:197] + "..."
			}
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

	return &Notification{
		Title: title,
		Body:  body,
	}
}

// formatNotificationEvent formats a Notification hook event
// Returns nil for permission-related notifications (handled by PermissionRequest hook)
func formatNotificationEvent(event *HookEvent) *Notification {
	// Skip permission-related notifications to avoid duplicates
	// These are already handled by the PermissionRequest hook
	if event.NotificationType == "permission_prompt" {
		return nil
	}

	body := event.Message
	if body == "" {
		body = "Notification from Claude Code"
	}

	return &Notification{
		Title: "Claude Code",
		Body:  body,
	}
}

// formatStop formats a Stop event
func formatStop(event *HookEvent) *Notification {
	body := "Task completed"
	if event.Cwd != "" {
		body = event.Cwd
	}

	return &Notification{
		Title: "Task Complete",
		Body:  body,
	}
}

// formatSessionStart formats a SessionStart event
func formatSessionStart(event *HookEvent) *Notification {
	body := "New session"
	if event.Cwd != "" {
		body = event.Cwd
	}

	return &Notification{
		Title: "Session Started",
		Body:  body,
	}
}
