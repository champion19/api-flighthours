package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type fakeMessageServiceReadOnly struct {
	getByIDRes    *domain.Message
	getByIDErr    error
	getByCodeRes  *domain.Message
	getByCodeErr  error
	listRes       []domain.Message
	listErr       error
	listActiveRes []domain.Message
	listActiveErr error
}

func (f *fakeMessageServiceReadOnly) BeginTx(context.Context) (output.Tx, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMessageServiceReadOnly) ValidateMessage(context.Context, domain.Message) error {
	return nil
}
func (f *fakeMessageServiceReadOnly) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	return f.getByIDRes, f.getByIDErr
}
func (f *fakeMessageServiceReadOnly) GetMessageByCode(ctx context.Context, code string) (*domain.Message, error) {
	return f.getByCodeRes, f.getByCodeErr
}
func (f *fakeMessageServiceReadOnly) ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error) {
	return f.listRes, f.listErr
}
func (f *fakeMessageServiceReadOnly) ListActiveMessages(ctx context.Context) ([]domain.Message, error) {
	return f.listActiveRes, f.listActiveErr
}
func (f *fakeMessageServiceReadOnly) SaveMessageToDB(context.Context, output.Tx, domain.Message) error {
	return nil
}
func (f *fakeMessageServiceReadOnly) UpdateMessageInDB(context.Context, output.Tx, domain.Message) error {
	return nil
}
func (f *fakeMessageServiceReadOnly) DeleteMessageFromDB(context.Context, output.Tx, string) error {
	return nil
}

var _ input.MessageService = (*fakeMessageServiceReadOnly)(nil)

func TestMessageInteractor_GetMessageByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		msg := &domain.Message{ID: "msg-123", Code: "TEST_001", Title: "Test"}
		svc := &fakeMessageServiceReadOnly{getByIDRes: msg}
		inter := NewMessageInteractor(svc, noopLogger{})

		result, err := inter.GetMessageByID(context.Background(), "msg-123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != "msg-123" {
			t.Fatalf("expected ID msg-123, got %s", result.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeMessageServiceReadOnly{getByIDErr: domain.ErrMessageNotFound}
		inter := NewMessageInteractor(svc, noopLogger{})

		_, err := inter.GetMessageByID(context.Background(), "nonexistent")

		if err != domain.ErrMessageNotFound {
			t.Fatalf("expected ErrMessageNotFound, got %v", err)
		}
	})
}

func TestMessageInteractor_GetMessageByCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		msg := &domain.Message{ID: "msg-123", Code: "TEST_001", Title: "Test"}
		svc := &fakeMessageServiceReadOnly{getByCodeRes: msg}
		inter := NewMessageInteractor(svc, noopLogger{})

		result, err := inter.GetMessageByCode(context.Background(), "TEST_001")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Code != "TEST_001" {
			t.Fatalf("expected code TEST_001, got %s", result.Code)
		}
	})
}

func TestMessageInteractor_ListMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		msgs := []domain.Message{
			{ID: "1", Code: "TEST_001", Title: "Test 1"},
			{ID: "2", Code: "TEST_002", Title: "Test 2"},
		}
		svc := &fakeMessageServiceReadOnly{listRes: msgs}
		inter := NewMessageInteractor(svc, noopLogger{})

		filters := map[string]interface{}{"module": "test"}
		result, err := inter.ListMessages(context.Background(), filters)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 messages, got %d", len(result))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		svc := &fakeMessageServiceReadOnly{listRes: []domain.Message{}}
		inter := NewMessageInteractor(svc, noopLogger{})

		result, err := inter.ListMessages(context.Background(), nil)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected 0 messages, got %d", len(result))
		}
	})
}

func TestMessageInteractor_ListActiveMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		msgs := []domain.Message{
			{ID: "1", Code: "TEST_001", Active: true},
			{ID: "2", Code: "TEST_002", Active: true},
		}
		svc := &fakeMessageServiceReadOnly{listActiveRes: msgs}
		inter := NewMessageInteractor(svc, noopLogger{})

		result, err := inter.ListActiveMessages(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 active messages, got %d", len(result))
		}
	})
}
