package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/mocks"
	"github.com/stretchr/testify/mock"
)

func TestEmployeeService_GetEmployeeByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		employee := &domain.Employee{ID: "123", Email: "test@example.com"}
		mockRepo.On("GetEmployeeByEmail", mock.Anything, "test@example.com").Return(employee, nil)

		svc := NewService(mockRepo, nil)
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

		svc := NewService(mockRepo, nil)
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

		svc := NewService(mockRepo, nil)
		result, err := svc.GetEmployeeByID(context.Background(), "123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != "123" {
			t.Fatalf("expected ID 123, got %s", result.ID)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockRepo.On("GetEmployeeByID", mock.Anything, "not-found").Return(nil, domain.ErrPersonNotFound)

		svc := NewService(mockRepo, nil)
		_, err := svc.GetEmployeeByID(context.Background(), "not-found")

		if !errors.Is(err, domain.ErrPersonNotFound) {
			t.Fatalf("expected ErrPersonNotFound, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_BeginTx(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockTx := new(mocks.MockTx)
	mockRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

	svc := NewService(mockRepo, nil)
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

		svc := NewService(mockRepo, nil)
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

		svc := NewService(mockRepo, nil)
		err := svc.SaveEmployeeToDB(context.Background(), mockTx, employee)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_UpdateEmployeeKeycloakID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockTx := new(mocks.MockTx)

		mockRepo.On("PatchEmployee", mock.Anything, mockTx, "emp123", "kc456").Return(nil)

		svc := NewService(mockRepo, nil)
		err := svc.UpdateEmployeeKeycloakID(context.Background(), mockTx, "emp123", "kc456")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockTx := new(mocks.MockTx)
		patchErr := errors.New("db error")

		mockRepo.On("PatchEmployee", mock.Anything, mockTx, "emp-fail", "kc-fail").Return(patchErr)

		svc := NewService(mockRepo, nil)
		err := svc.UpdateEmployeeKeycloakID(context.Background(), mockTx, "emp-fail", "kc-fail")

		if !errors.Is(err, patchErr) {
			t.Fatalf("expected %v, got %v", patchErr, err)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_SetUserPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SetPassword", mock.Anything, "user123", "newpass", false).Return(nil)

		svc := NewService(nil, mockAuth)
		err := svc.SetUserPassword(context.Background(), "user123", "newpass")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("keycloak error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SetPassword", mock.Anything, "user123", "newpass", false).Return(errors.New("keycloak error"))

		svc := NewService(nil, mockAuth)
		err := svc.SetUserPassword(context.Background(), "user123", "newpass")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_AssignUserRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("AssignRole", mock.Anything, "user123", "admin").Return(nil)

		svc := NewService(nil, mockAuth)
		err := svc.AssignUserRole(context.Background(), "user123", "admin")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("keycloak error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		roleErr := errors.New("keycloak error")
		mockAuth.On("AssignRole", mock.Anything, "user999", "admin").Return(roleErr)

		svc := NewService(nil, mockAuth)
		err := svc.AssignUserRole(context.Background(), "user999", "admin")

		if !errors.Is(err, roleErr) {
			t.Fatalf("expected %v, got %v", roleErr, err)
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_RollbackKeycloakUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("DeleteUser", mock.Anything, "kc123").Return(nil)

		svc := NewService(nil, mockAuth)
		err := svc.RollbackKeycloakUser(context.Background(), "kc123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("keycloak error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		deleteErr := errors.New("keycloak delete error")
		mockAuth.On("DeleteUser", mock.Anything, "kc-fail").Return(deleteErr)

		svc := NewService(nil, mockAuth)
		err := svc.RollbackKeycloakUser(context.Background(), "kc-fail")

		if !errors.Is(err, deleteErr) {
			t.Fatalf("expected %v, got %v", deleteErr, err)
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_LocateEmployee(t *testing.T) {
	t.Run("found in db", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		employee := &domain.Employee{ID: "123", Email: "test@example.com"}
		mockRepo.On("GetEmployeeByID", mock.Anything, "123").Return(employee, nil)

		svc := NewService(mockRepo, nil)
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

		svc := NewService(mockRepo, nil)
		_, err := svc.LocateEmployee(context.Background(), "999")

		if err != domain.ErrPersonNotFound {
			t.Fatalf("expected ErrEmployeeNotFound, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("employee nil", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		// Repo returns nil, nil (no error but no employee)
		mockRepo.On("GetEmployeeByID", mock.Anything, "empty").Return(nil, nil)

		svc := NewService(mockRepo, nil)
		_, err := svc.LocateEmployee(context.Background(), "empty")

		if err != domain.ErrPersonNotFound {
			t.Fatalf("expected ErrPersonNotFound, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_GetEmployeeByKeycloakID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		employee := &domain.Employee{ID: "123", KeycloakUserID: "kc-456"}
		mockRepo.On("GetEmployeeByKeycloakID", mock.Anything, "kc-456").Return(employee, nil)

		svc := NewService(mockRepo, nil)
		result, err := svc.GetEmployeeByKeycloakID(context.Background(), "kc-456")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.KeycloakUserID != "kc-456" {
			t.Fatalf("expected KeycloakUserID kc-456, got %s", result.KeycloakUserID)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockRepo.On("GetEmployeeByKeycloakID", mock.Anything, "kc-999").Return(nil, domain.ErrPersonNotFound)

		svc := NewService(mockRepo, nil)
		_, err := svc.GetEmployeeByKeycloakID(context.Background(), "kc-999")

		if err != domain.ErrPersonNotFound {
			t.Fatalf("expected ErrPersonNotFound, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockTx := new(mocks.MockTx)
		employee := domain.Employee{ID: "123", Email: "updated@example.com"}

		mockRepo.On("UpdateEmployee", mock.Anything, mockTx, employee).Return(nil)

		svc := NewService(mockRepo, nil)
		err := svc.UpdateEmployee(context.Background(), mockTx, employee)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockTx := new(mocks.MockTx)
		employee := domain.Employee{ID: "123", Email: "updated@example.com"}

		mockRepo.On("UpdateEmployee", mock.Anything, mockTx, employee).Return(errors.New("db error"))

		svc := NewService(mockRepo, nil)
		err := svc.UpdateEmployee(context.Background(), mockTx, employee)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {
	t.Run("success with keycloak user", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockAuth := new(mocks.MockAuthClient)

		mockAuth.On("DeleteUser", mock.Anything, "kc-123").Return(nil)
		mockRepo.On("DeleteEmployee", mock.Anything, nil, "emp-123").Return(nil)

		svc := NewService(mockRepo, mockAuth)
		err := svc.DeleteEmployee(context.Background(), "emp-123", "kc-123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("success without keycloak user", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)

		mockRepo.On("DeleteEmployee", mock.Anything, nil, "emp-123").Return(nil)

		svc := NewService(mockRepo, nil)
		err := svc.DeleteEmployee(context.Background(), "emp-123", "")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("keycloak delete error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("DeleteUser", mock.Anything, "kc-123").Return(errors.New("keycloak error"))

		svc := NewService(nil, mockAuth)
		err := svc.DeleteEmployee(context.Background(), "emp-123", "kc-123")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockAuth := new(mocks.MockAuthClient)

		mockAuth.On("DeleteUser", mock.Anything, "kc-123").Return(nil)
		mockRepo.On("DeleteEmployee", mock.Anything, nil, "emp-123").Return(errors.New("delete error"))

		svc := NewService(mockRepo, mockAuth)
		err := svc.DeleteEmployee(context.Background(), "emp-123", "kc-123")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_GetEmployeesByRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		employees := []domain.Employee{
			{ID: "1", Email: "admin1@example.com"},
			{ID: "2", Email: "admin2@example.com"},
		}
		mockRepo.On("GetEmployeesByRole", mock.Anything, "admin").Return(employees, nil)

		svc := NewService(mockRepo, nil)
		result, err := svc.GetEmployeesByRole(context.Background(), "admin")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 employees, got %d", len(result))
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty result", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockRepo.On("GetEmployeesByRole", mock.Anything, "unknown_role").Return([]domain.Employee{}, nil)

		svc := NewService(mockRepo, nil)
		result, err := svc.GetEmployeesByRole(context.Background(), "unknown_role")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected 0 employees, got %d", len(result))
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		repoErr := errors.New("db error")
		mockRepo.On("GetEmployeesByRole", mock.Anything, "admin").Return(nil, repoErr)

		svc := NewService(mockRepo, nil)
		_, err := svc.GetEmployeesByRole(context.Background(), "admin")

		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
		mockRepo.AssertExpectations(t)
	})
}
