//go:build integration
// +build integration

package message

import (
	"context"
	"testing"

	"github.com/champion19/api-flighthours/platform/databases/testhelper"
)

var testContainer *testhelper.MySQLContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testhelper.StartMySQL(ctx)
	if err != nil {
		panic("Failed to start MySQL container: " + err.Error())
	}
	testContainer = container

	if err := testContainer.SetupMessageSchema(ctx); err != nil {
		panic("Failed to setup schema: " + err.Error())
	}

	m.Run()

	testContainer.Stop(ctx)
}

func TestNewMessageRepository_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}

	t.Run("creates repository successfully", func(t *testing.T) {
		repo, err := NewMessageRepository(testContainer.DB)
		if err != nil {
			t.Fatalf("failed to create repository: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("fails with nil db", func(t *testing.T) {
		_, err := NewMessageRepository(nil)
		if err == nil {
			t.Error("expected error with nil db")
		}
	})
}

func TestRepository_GetMessageByCode_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanMessageTable(ctx)
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "msg-1", Code: "ERR001", Type: "error", Category: "validation", Module: "auth", Title: "Error Title", Content: "Error content", Active: true})

	repo, _ := NewMessageRepository(testContainer.DB)

	t.Run("finds existing message by code", func(t *testing.T) {
		msg, err := repo.GetByCode(ctx, "ERR001")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg == nil {
			t.Fatal("expected non-nil message")
		}
		if msg.Code != "ERR001" {
			t.Errorf("expected 'ERR001', got %q", msg.Code)
		}
	})

	t.Run("returns nil for non-existent code", func(t *testing.T) {
		msg, err := repo.GetByCode(ctx, "NONEXISTENT")
		// Repository returns nil, nil when not found
		if err != nil {
			t.Logf("info: returned error %v", err)
		}
		if msg != nil && msg.Code != "" {
			t.Error("expected nil or empty message for non-existent code")
		}
	})
}

func TestRepository_GetMessageByID_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanMessageTable(ctx)
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "msg-id-123", Code: "INFO001", Type: "info", Category: "general", Module: "system", Title: "Info Title", Content: "Info content", Active: true})

	repo, _ := NewMessageRepository(testContainer.DB)

	t.Run("finds existing message by ID", func(t *testing.T) {
		msg, err := repo.GetByID(ctx, "msg-id-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg == nil {
			t.Fatal("expected non-nil message")
		}
		if msg.ID != "msg-id-123" {
			t.Errorf("expected 'msg-id-123', got %q", msg.ID)
		}
	})

	t.Run("returns nil for non-existent ID", func(t *testing.T) {
		msg, err := repo.GetByID(ctx, "non-existent-id")
		// Repository returns nil, nil when not found
		if err != nil {
			t.Logf("info: returned error %v", err)
		}
		if msg != nil && msg.ID != "" {
			t.Error("expected nil or empty message for non-existent ID")
		}
	})
}

func TestRepository_ListActiveMessages_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanMessageTable(ctx)
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "m1", Code: "MSG001", Type: "info", Category: "cat1", Module: "mod1", Title: "Title1", Content: "Content1", Active: true})
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "m2", Code: "MSG002", Type: "error", Category: "cat2", Module: "mod2", Title: "Title2", Content: "Content2", Active: true})
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "m3", Code: "MSG003", Type: "warning", Category: "cat3", Module: "mod3", Title: "Title3", Content: "Content3", Active: false}) // inactive

	repo, _ := NewMessageRepository(testContainer.DB)

	t.Run("lists only active messages", func(t *testing.T) {
		messages, err := repo.GetAllActive(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(messages) != 2 {
			t.Errorf("expected 2 active messages, got %d", len(messages))
		}
	})
}

func TestRepository_GetByType_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	testContainer.CleanMessageTable(ctx)
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "e1", Code: "ERR001", Type: "error", Category: "cat", Module: "mod", Title: "Error1", Content: "Content", Active: true})
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "e2", Code: "ERR002", Type: "error", Category: "cat", Module: "mod", Title: "Error2", Content: "Content", Active: true})
	testContainer.InsertMessage(ctx, testhelper.MessageData{ID: "i1", Code: "INFO001", Type: "info", Category: "cat", Module: "mod", Title: "Info1", Content: "Content", Active: true})

	repo, _ := NewMessageRepository(testContainer.DB)

	t.Run("filters by type", func(t *testing.T) {
		messages, err := repo.GetByType(ctx, "error")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(messages) != 2 {
			t.Errorf("expected 2 error messages, got %d", len(messages))
		}
	})
}

func TestRepository_BeginTx_Integration(t *testing.T) {
	if testContainer == nil {
		t.Skip("Integration tests require docker")
	}
	ctx := context.Background()

	repo, _ := NewMessageRepository(testContainer.DB)

	t.Run("begins transaction successfully", func(t *testing.T) {
		tx, err := repo.BeginTx(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx == nil {
			t.Error("expected non-nil transaction")
		}
		tx.Rollback()
	})
}
