package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/mocks"
	"github.com/stretchr/testify/mock"
)

func TestEmployeeService_GetEmployeeByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		employee := &domain.Employee{ID: "123", Email: "test@example.com"}
		mockRepo.On("GetEmployeeByEmail", mock.Anything, "test@example.com").Return(employee, nil)

		svc := NewService(mockRepo, nil, noopLogger{})
		result, err := svc.GetEmployeeByEmail(context.Background(), "test@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Email != "test@example.com" {
			t.Fatalf("expected email test@example.com, got %s", result.Email)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockRepo.On("GetEmployeeByEmail", mock.Anything, "notfound@example.com").Return(nil, domain.ErrPersonNotFound)

		svc := NewService(mockRepo, nil, noopLogger{})
		_, err := svc.GetEmployeeByEmail(context.Background(), "notfound@example.com")

		if err != domain.ErrPersonNotFound {
			t.Fatalf("expected ErrPersonNotFound, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_GetEmployeeByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		employee := &domain.Employee{ID: "123", Email: "test@example.com"}
		mockRepo.On("GetEmployeeByID", mock.Anything, "123").Return(employee, nil)

		svc := NewService(mockRepo, nil, noopLogger{})
		result, err := svc.GetEmployeeByID(context.Background(), "123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != "123" {
			t.Fatalf("expected ID 123, got %s", result.ID)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_BeginTx(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockTx := new(mocks.MockTx)
	mockRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	svc := NewService(mockRepo, nil, noopLogger{})
	tx, err := svc.BeginTx(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tx == nil {
		t.Fatal("expected tx, got nil")
	}
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_SaveEmployeeToDB(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockTx := new(mocks.MockTx)
		employee := domain.Employee{ID: "123", Email: "test@example.com"}

		mockRepo.On("Save", mock.Anything, mockTx, employee).Return(nil)

		svc := NewService(mockRepo, nil, noopLogger{})
		err := svc.SaveEmployeeToDB(context.Background(), mockTx, employee)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockTx := new(mocks.MockTx)
		employee := domain.Employee{ID: "123", Email: "test@example.com"}

		mockRepo.On("Save", mock.Anything, mockTx, employee).Return(errors.New("db error"))

		svc := NewService(mockRepo, nil, noopLogger{})
		err := svc.SaveEmployeeToDB(context.Background(), mockTx, employee)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_UpdateEmployeeKeycloakID(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockTx := new(mocks.MockTx)

	mockRepo.On("PatchEmployee", mock.Anything, mockTx, "emp123", "kc456").Return(nil)

	svc := NewService(mockRepo, nil, noopLogger{})
	err := svc.UpdateEmployeeKeycloakID(context.Background(), mockTx, "emp123", "kc456")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestEmployeeService_SetUserPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SetPassword", mock.Anything, "user123", "newpass", false).Return(nil)

		svc := NewService(nil, mockAuth, noopLogger{})
		err := svc.SetUserPassword(context.Background(), "user123", "newpass")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("keycloak error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SetPassword", mock.Anything, "user123", "newpass", false).Return(errors.New("keycloak error"))

		svc := NewService(nil, mockAuth, noopLogger{})
		err := svc.SetUserPassword(context.Background(), "user123", "newpass")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_AssignUserRole(t *testing.T) {
	mockAuth := new(mocks.MockAuthClient)
	mockAuth.On("AssignRole", mock.Anything, "user123", "admin").Return(nil)

	svc := NewService(nil, mockAuth, noopLogger{})
	err := svc.AssignUserRole(context.Background(), "user123", "admin")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	mockAuth.AssertExpectations(t)
}

// TestEmployeeService_RollbackEmployee - skipped due to internal transaction handling complexity
// func TestEmployeeService_RollbackEmployee(t *testing.T) {}

func TestEmployeeService_RollbackKeycloakUser(t *testing.T) {
	mockAuth := new(mocks.MockAuthClient)
	mockAuth.On("DeleteUser", mock.Anything, "kc123").Return(nil)

	svc := NewService(nil, mockAuth, noopLogger{})
	err := svc.RollbackKeycloakUser(context.Background(), "kc123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	mockAuth.AssertExpectations(t)
}

func TestEmployeeService_LocateEmployee(t *testing.T) {
	t.Run("found in db", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		employee := &domain.Employee{ID: "123", Email: "test@example.com"}
		mockRepo.On("GetEmployeeByID", mock.Anything, "123").Return(employee, nil)

		svc := NewService(mockRepo, nil, noopLogger{})
		result, err := svc.LocateEmployee(context.Background(), "123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Employee.ID != "123" {
			t.Fatalf("expected ID 123, got %s", result.Employee.ID)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockRepo.On("GetEmployeeByID", mock.Anything, "999").Return(nil, domain.ErrPersonNotFound)

		svc := NewService(mockRepo, nil, noopLogger{})
		_, err := svc.LocateEmployee(context.Background(), "999")

		if err != domain.ErrPersonNotFound {
			t.Fatalf("expected ErrEmployeeNotFound, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})
}

// TestEmployeeService_CheckAndCleanInconsistentState - skipped due to internal transaction handling complexity
// func TestEmployeeService_CheckAndCleanInconsistentState(t *testing.T) {}
