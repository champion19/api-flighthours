package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type fakeMsgService struct {
	validateFn func(ctx context.Context, msg domain.Message) error
	beginTxFn  func(ctx context.Context) (output.Tx, error)
	saveFn     func(ctx context.Context, tx output.Tx, msg domain.Message) error
	getByIDFn  func(ctx context.Context, id string) (*domain.Message, error)
	updateFn   func(ctx context.Context, tx output.Tx, msg domain.Message) error
	deleteFn   func(ctx context.Context, tx output.Tx, id string) error
}

var _ input.MessageService = (*fakeMsgService)(nil)

func (f *fakeMsgService) BeginTx(ctx context.Context) (output.Tx, error) { return f.beginTxFn(ctx) }
func (f *fakeMsgService) ValidateMessage(ctx context.Context, msg domain.Message) error {
	return f.validateFn(ctx, msg)
}
func (f *fakeMsgService) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeMsgService) GetMessageByCode(context.Context, string) (*domain.Message, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMsgService) ListMessages(context.Context, map[string]interface{}) ([]domain.Message, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMsgService) ListActiveMessages(context.Context) ([]domain.Message, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMsgService) SaveMessageToDB(ctx context.Context, tx output.Tx, msg domain.Message) error {
	return f.saveFn(ctx, tx, msg)
}
func (f *fakeMsgService) UpdateMessageInDB(ctx context.Context, tx output.Tx, msg domain.Message) error {
	return f.updateFn(ctx, tx, msg)
}
func (f *fakeMsgService) DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error {
	return f.deleteFn(ctx, tx, id)
}

func TestMessageInteractor_CreateMessage(t *testing.T) {
	ctx := context.Background()
	msg := domain.Message{Code: "C1"}

	t.Run("validate fails => returns error, no tx", func(t *testing.T) {
		validateErr := errors.New("invalid")
		calledBegin := 0
		svc := &fakeMsgService{
			validateFn: func(context.Context, domain.Message) error { return validateErr },
			beginTxFn:  func(context.Context) (output.Tx, error) { calledBegin++; return nil, nil },
			saveFn:     func(context.Context, output.Tx, domain.Message) error { return nil },
		}
		i := NewMessageInteractor(svc, noopLogger{})

		_, err := i.CreateMessage(ctx, msg)
		if !errors.Is(err, validateErr) {
			t.Fatalf("expected %v, got %v", validateErr, err)
		}
		if calledBegin != 0 {
			t.Fatalf("expected BeginTx not called")
		}
	})

	t.Run("save fails => rollback", func(t *testing.T) {
		tx := &fakeTx{}
		saveErr := errors.New("save failed")
		svc := &fakeMsgService{
			validateFn: func(context.Context, domain.Message) error { return nil },
			beginTxFn:  func(context.Context) (output.Tx, error) { return tx, nil },
			saveFn:     func(context.Context, output.Tx, domain.Message) error { return saveErr },
		}
		i := NewMessageInteractor(svc, noopLogger{})

		_, err := i.CreateMessage(ctx, msg)
		if !errors.Is(err, saveErr) {
			t.Fatalf("expected %v, got %v", saveErr, err)
		}
		if !tx.rolledBack {
			t.Fatalf("expected rollback")
		}
		if tx.committed {
			t.Fatalf("did not expect commit")
		}
	})

	t.Run("commit fails => rollback", func(t *testing.T) {
		commitErr := errors.New("commit failed")
		tx := &fakeTx{commitFn: func() error { return commitErr }}
		svc := &fakeMsgService{
			validateFn: func(context.Context, domain.Message) error { return nil },
			beginTxFn:  func(context.Context) (output.Tx, error) { return tx, nil },
			saveFn:     func(context.Context, output.Tx, domain.Message) error { return nil },
		}
		i := NewMessageInteractor(svc, noopLogger{})

		_, err := i.CreateMessage(ctx, msg)
		if !errors.Is(err, commitErr) {
			t.Fatalf("expected error")
		}
		if !tx.committed {
			t.Fatalf("expected commit called")
		}
		if !tx.rolledBack {
			t.Fatalf("expected rollback on commit error")
		}
	})

	t.Run("success => commit true, rollback false", func(t *testing.T) {
		tx := &fakeTx{}
		svc := &fakeMsgService{
			validateFn: func(context.Context, domain.Message) error { return nil },
			beginTxFn:  func(context.Context) (output.Tx, error) { return tx, nil },
			saveFn:     func(context.Context, output.Tx, domain.Message) error { return nil },
		}
		i := NewMessageInteractor(svc, noopLogger{})

		res, err := i.CreateMessage(ctx, msg)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if res == nil {
			t.Fatalf("expected non-nil result")
		}
		if !tx.committed || tx.rolledBack {
			t.Fatalf("unexpected tx state committed=%v rolledBack=%v", tx.committed, tx.rolledBack)
		}
	})
}

func TestMessageInteractor_UpdateAndDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("update: not found => error", func(t *testing.T) {
		nf := domain.ErrMessageNotFound
		svc := &fakeMsgService{
			getByIDFn: func(context.Context, string) (*domain.Message, error) { return nil, nf },
			validateFn: func(context.Context, domain.Message) error {
				return nil
			},
			beginTxFn: func(context.Context) (output.Tx, error) { return &fakeTx{}, nil },
			updateFn:  func(context.Context, output.Tx, domain.Message) error { return nil },
			deleteFn:  func(context.Context, output.Tx, string) error { return nil },
			saveFn:    func(context.Context, output.Tx, domain.Message) error { return nil },
		}
		i := NewMessageInteractor(svc, noopLogger{})

		_, err := i.UpdateMessage(ctx, domain.Message{ID: "x", Code: "C"})
		if !errors.Is(err, nf) {
			t.Fatalf("expected %v, got %v", nf, err)
		}
	})

	t.Run("delete: fails in db => rollback", func(t *testing.T) {
		tx := &fakeTx{}
		deleteErr := errors.New("delete failed")
		svc := &fakeMsgService{
			getByIDFn: func(context.Context, string) (*domain.Message, error) { return &domain.Message{ID: "x"}, nil },
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			deleteFn:  func(context.Context, output.Tx, string) error { return deleteErr },
			validateFn: func(context.Context, domain.Message) error {
				return nil
			},
			updateFn: func(context.Context, output.Tx, domain.Message) error { return nil },
			saveFn:   func(context.Context, output.Tx, domain.Message) error { return nil },
		}
		i := NewMessageInteractor(svc, noopLogger{})

		err := i.DeleteMessage(ctx, "x")
		if !errors.Is(err, deleteErr) {
			t.Fatalf("expected %v, got %v", deleteErr, err)
		}
		if !tx.rolledBack {
			t.Fatalf("expected rollback")
		}
	})
}
