package services

import "context"

// Mock transaction for testing. Shared across this package's test files.
type mockTx struct {
	commitFn   func() error
	rollbackFn func() error
}

func (m *mockTx) Commit() error {
	if m.commitFn != nil {
		return m.commitFn()
	}
	return nil
}

func (m *mockTx) Rollback() error {
	if m.rollbackFn != nil {
		return m.rollbackFn()
	}
	return nil
}

func (m *mockTx) ExecContext(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (m *mockTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) interface{} {
	return nil
}
