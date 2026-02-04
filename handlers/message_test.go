package handlers

import (
	"errors"
	"testing"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestMessageRequest_Sanitize(t *testing.T) {
	t.Run("trims whitespace from all fields", func(t *testing.T) {
		msg := &MessageRequest{
			Code:     "  MSG001  ",
			Category: "  error  ",
			Module:   "  auth  ",
			Title:    "  Test Title  ",
			Content:  "  Some content  ",
		}

		msg.Sanitize()

		if msg.Code != "MSG001" {
			t.Errorf("expected Code 'MSG001', got %q", msg.Code)
		}
		if msg.Category != "error" {
			t.Errorf("expected Category 'error', got %q", msg.Category)
		}
		if msg.Module != "auth" {
			t.Errorf("expected Module 'auth', got %q", msg.Module)
		}
		if msg.Title != "Test Title" {
			t.Errorf("expected Title 'Test Title', got %q", msg.Title)
		}
		if msg.Content != "Some content" {
			t.Errorf("expected Content 'Some content', got %q", msg.Content)
		}
	})
}

func TestMessageRequest_ToDomain(t *testing.T) {
	t.Run("converts request to domain message", func(t *testing.T) {
		msg := MessageRequest{
			Code:     "MSG001",
			Type:     domain.MessageType("error"),
			Category: "validation",
			Module:   "core",
			Title:    "Error Title",
			Content:  "Error content here",
			Active:   true,
		}

		dm := msg.ToDomain()

		if dm.Code != "MSG001" {
			t.Errorf("expected Code 'MSG001', got %q", dm.Code)
		}
		if dm.Type != domain.MessageType("error") {
			t.Errorf("expected Type error, got %v", dm.Type)
		}
		if dm.Active != true {
			t.Error("expected Active to be true")
		}
	})
}

func TestToMessageResponse(t *testing.T) {
	t.Run("converts domain message to response", func(t *testing.T) {
		dm := &domain.Message{
			ID:       "msg-uuid-123",
			Code:     "INFO001",
			Type:     domain.MessageType("info"),
			Category: "system",
			Module:   "core",
			Title:    "Info Message",
			Content:  "Information content",
			Active:   true,
		}

		response := ToMessageResponse(dm, "encoded-123")

		if response.ID != "encoded-123" {
			t.Errorf("expected ID 'encoded-123', got %q", response.ID)
		}
		if response.Code != "INFO001" {
			t.Errorf("expected Code 'INFO001', got %q", response.Code)
		}
		if response.Active != true {
			t.Error("expected Active to be true")
		}
	})
}

func TestToMessageListResponse(t *testing.T) {
	t.Run("converts list of messages", func(t *testing.T) {
		messages := []domain.Message{
			{ID: "uuid-1", Code: "MSG001", Type: domain.MessageType("error"), Active: true},
			{ID: "uuid-2", Code: "MSG002", Type: domain.MessageType("info"), Active: false},
		}

		encodeFunc := func(id string) (string, error) {
			return "enc-" + id, nil
		}

		response := ToMessageListResponse(messages, encodeFunc)

		if len(response.Messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(response.Messages))
		}
		if response.Count != 2 {
			t.Errorf("expected count 2, got %d", response.Count)
		}
		if response.Messages[0].ID != "enc-uuid-1" {
			t.Errorf("expected encoded ID, got %q", response.Messages[0].ID)
		}
	})

	t.Run("handles encode errors gracefully", func(t *testing.T) {
		messages := []domain.Message{
			{ID: "uuid-1", Code: "MSG001"},
		}

		encodeFunc := func(id string) (string, error) {
			return "", errors.New("encode error")
		}

		response := ToMessageListResponse(messages, encodeFunc)

		// Should fall back to original ID on error
		if response.Messages[0].ID != "uuid-1" {
			t.Errorf("expected fallback to original ID, got %q", response.Messages[0].ID)
		}
	})

	t.Run("handles empty list", func(t *testing.T) {
		var messages []domain.Message

		response := ToMessageListResponse(messages, func(id string) (string, error) {
			return id, nil
		})

		if len(response.Messages) != 0 {
			t.Errorf("expected 0 messages, got %d", len(response.Messages))
		}
		if response.Count != 0 {
			t.Errorf("expected count 0, got %d", response.Count)
		}
	})
}

func TestCacheReloadResponse_Fields(t *testing.T) {
	t.Run("creates cache reload response", func(t *testing.T) {
		response := CacheReloadResponse{
			Success:     true,
			BeforeCount: 5,
			AfterCount:  10,
			Message:     "Cache reloaded",
		}

		if !response.Success {
			t.Error("expected Success to be true")
		}
		if response.BeforeCount != 5 {
			t.Errorf("expected BeforeCount 5, got %d", response.BeforeCount)
		}
		if response.AfterCount != 10 {
			t.Errorf("expected AfterCount 10, got %d", response.AfterCount)
		}
	})
}
