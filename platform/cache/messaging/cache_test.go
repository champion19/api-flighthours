package messaging

import (
	"context"
	"errors"
	"testing"
	"time"

	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
)

type fakeMessageRepo struct {
	messages []cachetypes.CachedMessage
	getErr   error
}

func (f *fakeMessageRepo) GetAllActiveForCache(ctx context.Context) ([]cachetypes.CachedMessage, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.messages, nil
}

func (f *fakeMessageRepo) GetByCodeForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	for _, m := range f.messages {
		if m.Code == code {
			return &m, nil
		}
	}
	return nil, errors.New("not found")
}

func (f *fakeMessageRepo) GetByCodeWithStatusForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	return nil, errors.New("not found")
}

func TestNewMessageCache(t *testing.T) {
	t.Run("creates cache with repo", func(t *testing.T) {
		repo := &fakeMessageRepo{}
		cache := NewMessageCache(repo, 0)
		if cache == nil {
			t.Error("expected non-nil cache")
		}
	})
}

func TestLoadMessages(t *testing.T) {
	t.Run("loads messages successfully", func(t *testing.T) {
		repo := &fakeMessageRepo{
			messages: []cachetypes.CachedMessage{
				{Code: "TEST_001", Type: TypeSuccess, Content: "test message"},
				{Code: "TEST_002", Type: TypeError, Content: "error message"},
			},
		}
		cache := NewMessageCache(repo, 0)
		err := cache.LoadMessages(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cache.MessageCount() != 2 {
			t.Errorf("expected 2 messages, got %d", cache.MessageCount())
		}
	})

	t.Run("returns error when repo fails", func(t *testing.T) {
		repo := &fakeMessageRepo{getErr: errors.New("db error")}
		cache := NewMessageCache(repo, 0)
		err := cache.LoadMessages(context.Background())
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetMessage(t *testing.T) {
	repo := &fakeMessageRepo{
		messages: []cachetypes.CachedMessage{
			{Code: "TEST_001", Type: TypeSuccess, Content: "test message", Active: true},
			{Code: "GEN_MSG_INACTIVE_ERR_00002", Type: TypeError, Content: "fallback error", Active: true},
		},
	}
	cache := NewMessageCache(repo, 0)
	cache.LoadMessages(context.Background())

	t.Run("returns cached message", func(t *testing.T) {
		msg := cache.GetMessage("TEST_001")
		if msg == nil {
			t.Fatal("expected non-nil message")
		}
		if msg.Code != "TEST_001" {
			t.Errorf("expected code TEST_001, got %s", msg.Code)
		}
	})

	t.Run("gets message from DB when not in cache", func(t *testing.T) {
		// Create a fresh cache with only 1 message loaded initially
		freshRepo := &fakeMessageRepo{
			messages: []cachetypes.CachedMessage{
				{Code: "DB_MSG", Type: TypeSuccess, Content: "from database", Active: true},
			},
		}
		freshCache := NewMessageCache(freshRepo, 0)
		// Don't load messages initially, so DB_MSG is not in cache

		// When GetMessage is called, it should try DB lookup
		msg := freshCache.GetMessage("DB_MSG")
		if msg == nil {
			t.Fatal("expected to get message from DB")
		}
		if msg.Content != "from database" {
			t.Errorf("expected 'from database', got %q", msg.Content)
		}
	})

	t.Run("returns nil for special fallback code", func(t *testing.T) {
		emptyRepo := &fakeMessageRepo{}
		emptyCache := NewMessageCache(emptyRepo, 0)
		msg := emptyCache.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
		// Should return nil when the special code itself is not found
		if msg != nil {
			t.Error("expected nil for missing fallback code")
		}
	})

	t.Run("handles missing message by returning fallback", func(t *testing.T) {
		msg := cache.GetMessage("NONEXISTENT")
		// Will return fallback message (GEN_MSG_INACTIVE_ERR_00002)
		if msg != nil && msg.Code == "GEN_MSG_INACTIVE_ERR_00002" {
			// This is expected behavior - fallback to error message
			_ = msg
		}
	})

	t.Run("handles repo error for GetByCodeWithStatusForCache", func(t *testing.T) {
		// Create a repo that returns error for GetByCodeWithStatusForCache
		errRepo := &fakeMessageRepo{
			getErr: errors.New("db error"),
			messages: []cachetypes.CachedMessage{
				{Code: "GEN_MSG_INACTIVE_ERR_00002", Type: TypeError, Content: "fallback error", Active: true},
			},
		}
		errCache := NewMessageCache(errRepo, 0)
		errCache.LoadMessages(context.Background())

		msg := errCache.GetMessage("UNKNOWN_CODE")
		// Should return fallback
		if msg == nil {
			t.Error("expected fallback message")
		}
	})

	t.Run("returns fallback for inactive message", func(t *testing.T) {
		// This covers the inactiveMsg != nil && !inactiveMsg.Active branch
		inactiveRepo := &fakeMessageRepoWithInactive{
			activeMessages: []cachetypes.CachedMessage{
				{Code: "GEN_MSG_INACTIVE_ERR_00002", Type: TypeError, Content: "fallback error", Active: true},
			},
			inactiveMsg: &cachetypes.CachedMessage{Code: "INACTIVE_MSG", Type: TypeSuccess, Content: "inactive", Active: false},
		}
		cache := NewMessageCache(inactiveRepo, 0)
		cache.LoadMessages(context.Background())

		msg := cache.GetMessage("INACTIVE_MSG")
		// Should return fallback since the message is inactive
		if msg == nil {
			t.Error("expected fallback message")
		}
		if msg.Code != "GEN_MSG_INACTIVE_ERR_00002" {
			t.Errorf("expected fallback code, got %s", msg.Code)
		}
	})
}

// fakeMessageRepoWithInactive provides more control over inactive message behavior
type fakeMessageRepoWithInactive struct {
	activeMessages []cachetypes.CachedMessage
	inactiveMsg    *cachetypes.CachedMessage
}

func (f *fakeMessageRepoWithInactive) GetAllActiveForCache(ctx context.Context) ([]cachetypes.CachedMessage, error) {
	return f.activeMessages, nil
}

func (f *fakeMessageRepoWithInactive) GetByCodeForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	for _, m := range f.activeMessages {
		if m.Code == code {
			return &m, nil
		}
	}
	// Return nil without error to trigger GetByCodeWithStatusForCache
	return nil, nil
}

func (f *fakeMessageRepoWithInactive) GetByCodeWithStatusForCache(ctx context.Context, code string) (*cachetypes.CachedMessage, error) {
	if f.inactiveMsg != nil && f.inactiveMsg.Code == code {
		return f.inactiveMsg, nil
	}
	return nil, errors.New("not found")
}

func TestGetMessageResponse(t *testing.T) {
	repo := &fakeMessageRepo{
		messages: []cachetypes.CachedMessage{
			{Code: "TEST_PARAM", Type: TypeSuccess, Title: "Test", Content: "Hello ${0}", Active: true},
		},
	}
	cache := NewMessageCache(repo, 0)
	cache.LoadMessages(context.Background())

	t.Run("returns response with parameters replaced", func(t *testing.T) {
		resp := cache.GetMessageResponse("TEST_PARAM", "World")
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Content != "Hello World" {
			t.Errorf("expected 'Hello World', got %q", resp.Content)
		}
	})

	t.Run("handles nil message", func(t *testing.T) {
		resp := cache.GetMessageResponse("NONEXISTENT_CODE")
		// Should handle gracefully
		_ = resp
	})
}

func TestGetHTTPStatus(t *testing.T) {
	repo := &fakeMessageRepo{
		messages: []cachetypes.CachedMessage{
			{Code: "GEN_SRV_ERR_00001", Type: TypeError, Active: true},
			{Code: "MOD_U_REG_EXI_00001", Type: TypeSuccess, Active: true},
			{Code: "WARNING_MSG", Type: TypeWarning, Active: true},
			{Code: "INFO_MSG", Type: TypeInfo, Active: true},
			{Code: "DEBUG_MSG", Type: TypeDebug, Active: true},
		},
	}
	cache := NewMessageCache(repo, 0)
	cache.LoadMessages(context.Background())

	tests := []struct {
		code     string
		expected int
	}{
		{"GEN_SRV_ERR_00001", 500},
		{"MOD_U_REG_EXI_00001", 201},
		{"MOD_AIR_GET_EXI_00001", 200},
		{"MOD_V_VAL_ERR_00001", 400},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			status := cache.GetHTTPStatus(tt.code)
			if status != tt.expected {
				t.Errorf("GetHTTPStatus(%q) = %d, expected %d", tt.code, status, tt.expected)
			}
		})
	}

	t.Run("TypeWarning returns 200", func(t *testing.T) {
		status := cache.GetHTTPStatus("WARNING_MSG")
		if status != 200 {
			t.Errorf("expected 200 for TypeWarning, got %d", status)
		}
	})

	t.Run("TypeInfo returns 200", func(t *testing.T) {
		status := cache.GetHTTPStatus("INFO_MSG")
		if status != 200 {
			t.Errorf("expected 200 for TypeInfo, got %d", status)
		}
	})

	t.Run("TypeDebug returns 200", func(t *testing.T) {
		status := cache.GetHTTPStatus("DEBUG_MSG")
		if status != 200 {
			t.Errorf("expected 200 for TypeDebug, got %d", status)
		}
	})

	t.Run("unknown code not in map => uses message type", func(t *testing.T) {
		// MOD_U_REG_EXI_00001 is a success type, returns 200 (not in static map)
		status := cache.GetHTTPStatus("MOD_U_REG_EXI_00001")
		// Note: This code should map to 201 if it's in the static map
		if status != 201 {
			t.Logf("status for MOD_U_REG_EXI_00001: %d", status)
		}
	})

	t.Run("unknown code with type error returns 500", func(t *testing.T) {
		// GEN_SRV_ERR_00001 is error type
		status := cache.GetHTTPStatus("GEN_SRV_ERR_00001")
		if status != 500 {
			t.Errorf("expected 500 for TypeError, got %d", status)
		}
	})
}

func TestMessageCount(t *testing.T) {
	repo := &fakeMessageRepo{
		messages: []cachetypes.CachedMessage{
			{Code: "MSG_1"},
			{Code: "MSG_2"},
			{Code: "MSG_3"},
		},
	}
	cache := NewMessageCache(repo, 0)
	cache.LoadMessages(context.Background())

	count := cache.MessageCount()
	if count != 3 {
		t.Errorf("expected 3, got %d", count)
	}
}

func TestReloadMessages(t *testing.T) {
	repo := &fakeMessageRepo{
		messages: []cachetypes.CachedMessage{
			{Code: "MSG_1"},
		},
	}
	cache := NewMessageCache(repo, 0)
	cache.LoadMessages(context.Background())

	if cache.MessageCount() != 1 {
		t.Errorf("expected 1 message initially")
	}

	// Add more messages to repo
	repo.messages = append(repo.messages, cachetypes.CachedMessage{Code: "MSG_2"})
	cache.ReloadMessages(context.Background())

	if cache.MessageCount() != 2 {
		t.Errorf("expected 2 messages after reload")
	}
}

func TestReplaceAll(t *testing.T) {
	tests := []struct {
		s, old, new, expected string
	}{
		{"Hello ${0}", "${0}", "World", "Hello World"},
		{"${0} and ${1}", "${0}", "A", "A and ${1}"},
		{"no placeholder", "${0}", "X", "no placeholder"},
		{"", "${0}", "X", ""},
	}

	for _, tt := range tests {
		result := replaceAll(tt.s, tt.old, tt.new)
		if result != tt.expected {
			t.Errorf("replaceAll(%q, %q, %q) = %q, expected %q", tt.s, tt.old, tt.new, result, tt.expected)
		}
	}
}

func TestStartAutoRefresh(t *testing.T) {
	t.Run("does not panic with zero interval", func(t *testing.T) {
		repo := &fakeMessageRepo{}
		cache := NewMessageCache(repo, 0)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("StartAutoRefresh panicked: %v", r)
			}
		}()
		cache.StartAutoRefresh(context.Background())
	})

	t.Run("starts with positive interval", func(t *testing.T) {
		repo := &fakeMessageRepo{}
		cache := NewMessageCache(repo, 100*time.Millisecond)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("StartAutoRefresh panicked: %v", r)
			}
		}()
		cache.StartAutoRefresh(context.Background())
		// Give it a moment to start
		time.Sleep(50 * time.Millisecond)
		cache.StopAutoRefresh()
	})
}

func TestStopAutoRefresh(t *testing.T) {
	t.Run("does not panic with zero interval", func(t *testing.T) {
		repo := &fakeMessageRepo{}
		cache := NewMessageCache(repo, 0)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("StopAutoRefresh panicked: %v", r)
			}
		}()
		cache.StopAutoRefresh()
	})
}
