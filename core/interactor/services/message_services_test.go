package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type fakeMsgRepo struct {
	getByCodeFn   func(ctx context.Context, code string) (*domain.Message, error)
	getByIDFn     func(ctx context.Context, id string) (*domain.Message, error)
	getAllActive  func(ctx context.Context) ([]domain.Message, error)
	saveMsgFn     func(ctx context.Context, tx output.Tx, message domain.Message) error
	updateMsgFn   func(ctx context.Context, tx output.Tx, message domain.Message) error
	deleteMsgFn   func(ctx context.Context, tx output.Tx, id string) error
	beginTxFn     func(ctx context.Context) (output.Tx, error)
	getByModuleFn func(ctx context.Context, module string) ([]domain.Message, error)
	getByTypeFn   func(ctx context.Context, msgType string) ([]domain.Message, error)
}

func (f fakeMsgRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) SaveMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	if f.saveMsgFn != nil {
		return f.saveMsgFn(ctx, tx, message)
	}
	return errors.New("not implemented")
}
func (f fakeMsgRepo) UpdateMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	if f.updateMsgFn != nil {
		return f.updateMsgFn(ctx, tx, message)
	}
	return errors.New("not implemented")
}
func (f fakeMsgRepo) DeleteMessage(ctx context.Context, tx output.Tx, id string) error {
	if f.deleteMsgFn != nil {
		return f.deleteMsgFn(ctx, tx, id)
	}
	return errors.New("not implemented")
}
func (f fakeMsgRepo) GetAllActive(ctx context.Context) ([]domain.Message, error) {
	if f.getAllActive != nil {
		return f.getAllActive(ctx)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByCode(ctx context.Context, code string) (*domain.Message, error) {
	if f.getByCodeFn != nil {
		return f.getByCodeFn(ctx, code)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByType(ctx context.Context, msgType string) ([]domain.Message, error) {
	if f.getByTypeFn != nil {
		return f.getByTypeFn(ctx, msgType)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	if f.getByModuleFn != nil {
		return f.getByModuleFn(ctx, module)
	}
	return nil, errors.New("not implemented")
}

type fakeTx struct{}

func (f *fakeTx) Commit() error   { return nil }
func (f *fakeTx) Rollback() error { return nil }

func TestMessageService_ValidateMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("empty code => ErrMessageCodeRequired", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{}, noopLogger{})
		err := svc.ValidateMessage(ctx, domain.Message{Code: ""})
		if !errors.Is(err, domain.ErrMessageCodeRequired) {
			t.Fatalf("expected %v, got %v", domain.ErrMessageCodeRequired, err)
		}
	})

	t.Run("create and code exists => ErrMessageCodeDuplicate", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{getByCodeFn: func(context.Context, string) (*domain.Message, error) {
			return &domain.Message{ID: "1", Code: "C1"}, nil
		}}, noopLogger{})

		err := svc.ValidateMessage(ctx, domain.Message{ID: "", Code: "C1"})
		if !errors.Is(err, domain.ErrMessageCodeDuplicate) {
			t.Fatalf("expected %v, got %v", domain.ErrMessageCodeDuplicate, err)
		}
	})

	t.Run("update (ID not empty) => does not check duplicate", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{getByCodeFn: func(context.Context, string) (*domain.Message, error) {
			return &domain.Message{ID: "1", Code: "C1"}, nil
		}}, noopLogger{})

		err := svc.ValidateMessage(ctx, domain.Message{ID: "existing", Code: "C1"})
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
}

func TestMessageService_GetMessageByID(t *testing.T) {
	ctx := context.Background()

	t.Run("repo returns nil,nil => ErrMessageNotFound", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{getByIDFn: func(context.Context, string) (*domain.Message, error) {
			return nil, nil
		}}, noopLogger{})

		_, err := svc.GetMessageByID(ctx, "x")
		if !errors.Is(err, domain.ErrMessageNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrMessageNotFound, err)
		}
	})

	t.Run("repo returns error => propagate", func(t *testing.T) {
		repoErr := errors.New("db error")
		svc := NewMessageService(fakeMsgRepo{getByIDFn: func(context.Context, string) (*domain.Message, error) {
			return nil, repoErr
		}}, noopLogger{})

		_, err := svc.GetMessageByID(ctx, "x")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		msg := &domain.Message{ID: "123", Code: "TEST_001"}
		svc := NewMessageService(fakeMsgRepo{getByIDFn: func(context.Context, string) (*domain.Message, error) {
			return msg, nil
		}}, noopLogger{})

		result, err := svc.GetMessageByID(ctx, "123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != "123" {
			t.Fatalf("expected ID 123, got %s", result.ID)
		}
	})
}

func TestMessageService_GetMessageByCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		msg := &domain.Message{ID: "123", Code: "TEST_001"}
		svc := NewMessageService(fakeMsgRepo{getByCodeFn: func(context.Context, string) (*domain.Message, error) {
			return msg, nil
		}}, noopLogger{})

		result, err := svc.GetMessageByCode(ctx, "TEST_001")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Code != "TEST_001" {
			t.Fatalf("expected code TEST_001, got %s", result.Code)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		svc := NewMessageService(fakeMsgRepo{getByCodeFn: func(context.Context, string) (*domain.Message, error) {
			return nil, repoErr
		}}, noopLogger{})

		_, err := svc.GetMessageByCode(ctx, "TEST_001")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestMessageService_ListMessages(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		msgs := []domain.Message{
			{ID: "1", Code: "TEST_001"},
			{ID: "2", Code: "TEST_002"},
		}
		svc := NewMessageService(fakeMsgRepo{getAllActive: func(context.Context) ([]domain.Message, error) {
			return msgs, nil
		}}, noopLogger{})

		result, err := svc.ListMessages(ctx, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 messages, got %d", len(result))
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		svc := NewMessageService(fakeMsgRepo{getAllActive: func(context.Context) ([]domain.Message, error) {
			return nil, repoErr
		}}, noopLogger{})

		_, err := svc.ListMessages(ctx, nil)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestMessageService_ListActiveMessages(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		msgs := []domain.Message{
			{ID: "1", Code: "TEST_001", Active: true},
			{ID: "2", Code: "TEST_002", Active: true},
		}
		svc := NewMessageService(fakeMsgRepo{getAllActive: func(context.Context) ([]domain.Message, error) {
			return msgs, nil
		}}, noopLogger{})

		result, err := svc.ListActiveMessages(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 messages, got %d", len(result))
		}
	})
}

func TestMessageService_SaveMessageToDB(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		msg := domain.Message{ID: "123", Code: "TEST_001"}
		svc := NewMessageService(fakeMsgRepo{saveMsgFn: func(context.Context, output.Tx, domain.Message) error {
			return nil
		}}, noopLogger{})

		err := svc.SaveMessageToDB(ctx, nil, msg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		msg := domain.Message{ID: "123", Code: "TEST_001"}
		svc := NewMessageService(fakeMsgRepo{saveMsgFn: func(context.Context, output.Tx, domain.Message) error {
			return repoErr
		}}, noopLogger{})

		err := svc.SaveMessageToDB(ctx, nil, msg)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestMessageService_UpdateMessageInDB(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		msg := domain.Message{ID: "123", Code: "TEST_001"}
		svc := NewMessageService(fakeMsgRepo{
			getByIDFn: func(context.Context, string) (*domain.Message, error) {
				return &domain.Message{ID: "123", Code: "TEST_001"}, nil
			},
			updateMsgFn: func(context.Context, output.Tx, domain.Message) error {
				return nil
			},
		}, noopLogger{})

		err := svc.UpdateMessageInDB(ctx, nil, msg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		msg := domain.Message{ID: "123", Code: "TEST_001"}
		svc := NewMessageService(fakeMsgRepo{
			getByIDFn: func(context.Context, string) (*domain.Message, error) {
				return &domain.Message{ID: "123", Code: "TEST_001"}, nil
			},
			updateMsgFn: func(context.Context, output.Tx, domain.Message) error {
				return repoErr
			},
		}, noopLogger{})

		err := svc.UpdateMessageInDB(ctx, nil, msg)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestMessageService_DeleteMessageFromDB(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{deleteMsgFn: func(context.Context, output.Tx, string) error {
			return nil
		}}, noopLogger{})

		err := svc.DeleteMessageFromDB(ctx, nil, "123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		svc := NewMessageService(fakeMsgRepo{deleteMsgFn: func(context.Context, output.Tx, string) error {
			return repoErr
		}}, noopLogger{})

		err := svc.DeleteMessageFromDB(ctx, nil, "123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestMessageService_BeginTx(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		fakeTx := &fakeTx{}
		svc := NewMessageService(fakeMsgRepo{beginTxFn: func(context.Context) (output.Tx, error) {
			return fakeTx, nil
		}}, noopLogger{})

		tx, err := svc.BeginTx(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tx == nil {
			t.Fatal("expected tx, got nil")
		}
	})
}
