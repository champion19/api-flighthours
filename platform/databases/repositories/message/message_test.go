package message

import (
	"testing"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestMessage_ToDomain(t *testing.T) {
	t.Run("converts Message to domain.Message", func(t *testing.T) {
		now := time.Now()
		msg := Message{
			ID:        "msg-123",
			Code:      "ERR001",
			Type:      "error",
			Category:  "validation",
			Module:    "auth",
			Title:     "Error Title",
			Content:   "Error content",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		result := msg.ToDomain()

		if result.ID != "msg-123" {
			t.Errorf("expected ID 'msg-123', got %q", result.ID)
		}
		if result.Code != "ERR001" {
			t.Errorf("expected Code 'ERR001', got %q", result.Code)
		}
		if result.Type != domain.MessageType("error") {
			t.Errorf("expected Type 'error', got %q", result.Type)
		}
		if result.Category != "validation" {
			t.Errorf("expected Category 'validation', got %q", result.Category)
		}
		if result.Module != "auth" {
			t.Errorf("expected Module 'auth', got %q", result.Module)
		}
		if result.Title != "Error Title" {
			t.Errorf("expected Title 'Error Title', got %q", result.Title)
		}
		if result.Content != "Error content" {
			t.Errorf("expected Content 'Error content', got %q", result.Content)
		}
		if !result.Active {
			t.Error("expected Active to be true")
		}
	})

	t.Run("converts inactive Message", func(t *testing.T) {
		msg := Message{
			ID:     "msg-456",
			Active: false,
		}

		result := msg.ToDomain()

		if result.Active {
			t.Error("expected Active to be false")
		}
	})
}

func TestFromDomain(t *testing.T) {
	t.Run("converts domain.Message to Message", func(t *testing.T) {
		now := time.Now()
		domainMsg := domain.Message{
			ID:        "msg-123",
			Code:      "INFO001",
			Type:      domain.MessageType("info"),
			Category:  "notification",
			Module:    "system",
			Title:     "Info Title",
			Content:   "Info content",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		result := FromDomain(domainMsg)

		if result.ID != "msg-123" {
			t.Errorf("expected ID 'msg-123', got %q", result.ID)
		}
		if result.Code != "INFO001" {
			t.Errorf("expected Code 'INFO001', got %q", result.Code)
		}
		if result.Type != "info" {
			t.Errorf("expected Type 'info', got %q", result.Type)
		}
		if result.Category != "notification" {
			t.Errorf("expected Category 'notification', got %q", result.Category)
		}
		if result.Module != "system" {
			t.Errorf("expected Module 'system', got %q", result.Module)
		}
		if result.Title != "Info Title" {
			t.Errorf("expected Title 'Info Title', got %q", result.Title)
		}
		if result.Content != "Info content" {
			t.Errorf("expected Content 'Info content', got %q", result.Content)
		}
		if !result.Active {
			t.Error("expected Active to be true")
		}
	})

	t.Run("preserves timestamps", func(t *testing.T) {
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

		domainMsg := domain.Message{
			ID:        "msg-789",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		result := FromDomain(domainMsg)

		if !result.CreatedAt.Equal(createdAt) {
			t.Error("CreatedAt should be preserved")
		}
		if !result.UpdatedAt.Equal(updatedAt) {
			t.Error("UpdatedAt should be preserved")
		}
	})
}

func TestMessageStruct(t *testing.T) {
	t.Run("creates Message with all fields", func(t *testing.T) {
		now := time.Now()
		msg := Message{
			ID:        "msg-123",
			Code:      "TEST001",
			Type:      "warning",
			Category:  "test",
			Module:    "testing",
			Title:     "Test Title",
			Content:   "Test content",
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if msg.ID != "msg-123" {
			t.Errorf("expected ID 'msg-123', got %q", msg.ID)
		}
		if msg.Type != "warning" {
			t.Errorf("expected Type 'warning', got %q", msg.Type)
		}
	})
}
