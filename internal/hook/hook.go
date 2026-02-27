package hook

import (
	"encoding/json"
	"fmt"
)

type Input struct {
	SessionID        string `json:"session_id"`
	HookEventName    string `json:"hook_event_name"`
	Message          string `json:"message,omitempty"`
	Title            string `json:"title,omitempty"`
	NotificationType string `json:"notification_type,omitempty"`
	StopHookActive   bool   `json:"stop_hook_active,omitempty"`
}

func Parse(data []byte) (*Input, error) {
	var input Input
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("failed to parse hook input: %w", err)
	}
	if input.HookEventName == "" {
		return nil, fmt.Errorf("missing hook_event_name")
	}
	return &input, nil
}

func (i *Input) FormatNotification() (title, message string) {
	switch i.HookEventName {
	case "Notification":
		title = "Claude Code"
		if i.Title != "" {
			title = i.Title
		}
		message = i.Message
		if message == "" {
			message = "New notification"
		}
	case "Stop":
		title = "Claude Code"
		message = "Claude has finished responding"
	default:
		title = "Claude Code"
		message = fmt.Sprintf("Event: %s", i.HookEventName)
	}
	return title, message
}
