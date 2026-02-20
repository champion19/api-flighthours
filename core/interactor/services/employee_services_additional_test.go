package services

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/Nerzal/gocloak/v13"
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

func TestEmployeeService_RollbackEmployee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockRepo.On("DeleteEmployee", mock.Anything, nil, "emp-123").Return(nil)

		svc := NewService(mockRepo, nil)
		err := svc.RollbackEmployee(context.Background(), "emp-123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		deleteErr := errors.New("delete error")
		mockRepo.On("DeleteEmployee", mock.Anything, nil, "emp-fail").Return(deleteErr)

		svc := NewService(mockRepo, nil)
		err := svc.RollbackEmployee(context.Background(), "emp-fail")

		if !errors.Is(err, deleteErr) {
			t.Fatalf("expected %v, got %v", deleteErr, err)
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{"exact match", "connection refused", "connection refused", true},
		{"case insensitive", "Connection Refused", "connection refused", true},
		{"substring", "error: connection refused by host", "connection refused", true},
		{"no match", "timeout error", "connection refused", false},
		{"empty substr", "anything", "", true},
		{"empty both", "", "", true},
		{"empty s", "", "something", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := contains(tc.s, tc.substr); got != tc.want {
				t.Errorf("contains(%q, %q) = %v, want %v", tc.s, tc.substr, got, tc.want)
			}
		})
	}
}

func TestIsConnectionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"connection refused", errors.New("connection refused"), true},
		{"no such host", errors.New("no such host"), true},
		{"connection reset", errors.New("connection reset"), true},
		{"network is unreachable", errors.New("network is unreachable"), true},
		{"unrelated error", errors.New("invalid input"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isConnectionError(tc.err); got != tc.want {
				t.Errorf("isConnectionError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"timeout", errors.New("timeout"), true},
		{"deadline exceeded", errors.New("deadline exceeded"), true},
		{"context deadline exceeded", errors.New("context deadline exceeded"), true},
		{"unrelated error", errors.New("invalid input"), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTimeoutError(tc.err); got != tc.want {
				t.Errorf("isTimeoutError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func strPtr(s string) *string { return &s }

func TestEmployeeService_CheckAndCleanInconsistentState(t *testing.T) {
	t.Run("both exist - no cleanup", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockAuth := new(mocks.MockAuthClient)

		employee := &domain.Employee{ID: "emp-123", Email: "test@example.com"}
		kcUser := &gocloak.User{ID: strPtr("kc-123"), Email: strPtr("test@example.com")}

		mockRepo.On("GetEmployeeByEmail", mock.Anything, "test@example.com").Return(employee, nil)
		mockAuth.On("GetUserByEmail", mock.Anything, "test@example.com").Return(kcUser, nil)

		svc := NewService(mockRepo, mockAuth)
		err := svc.CheckAndCleanInconsistentState(context.Background(), "test@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("neither exist - no cleanup", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockAuth := new(mocks.MockAuthClient)

		mockRepo.On("GetEmployeeByEmail", mock.Anything, "new@example.com").Return(nil, domain.ErrPersonNotFound)
		mockAuth.On("GetUserByEmail", mock.Anything, "new@example.com").Return(nil, errors.New("not found"))

		svc := NewService(mockRepo, mockAuth)
		err := svc.CheckAndCleanInconsistentState(context.Background(), "new@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("only in keycloak - cleanup", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockAuth := new(mocks.MockAuthClient)

		kcUser := &gocloak.User{ID: strPtr("kc-orphan"), Email: strPtr("orphan@example.com")}
		mockRepo.On("GetEmployeeByEmail", mock.Anything, "orphan@example.com").Return(nil, domain.ErrPersonNotFound)
		mockAuth.On("GetUserByEmail", mock.Anything, "orphan@example.com").Return(kcUser, nil)
		mockAuth.On("DeleteUser", mock.Anything, "kc-orphan").Return(nil)

		svc := NewService(mockRepo, mockAuth)
		err := svc.CheckAndCleanInconsistentState(context.Background(), "orphan@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("only in db - cleanup", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockAuth := new(mocks.MockAuthClient)

		employee := &domain.Employee{ID: "emp-orphan", Email: "dborphan@example.com"}
		mockRepo.On("GetEmployeeByEmail", mock.Anything, "dborphan@example.com").Return(employee, nil)
		mockAuth.On("GetUserByEmail", mock.Anything, "dborphan@example.com").Return(nil, errors.New("not found"))
		mockRepo.On("DeleteEmployee", mock.Anything, nil, "emp-orphan").Return(nil)

		svc := NewService(mockRepo, mockAuth)
		err := svc.CheckAndCleanInconsistentState(context.Background(), "dborphan@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("keycloak cleanup error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockAuth := new(mocks.MockAuthClient)

		kcUser := &gocloak.User{ID: strPtr("kc-fail"), Email: strPtr("fail@example.com")}
		mockRepo.On("GetEmployeeByEmail", mock.Anything, "fail@example.com").Return(nil, domain.ErrPersonNotFound)
		mockAuth.On("GetUserByEmail", mock.Anything, "fail@example.com").Return(kcUser, nil)
		mockAuth.On("DeleteUser", mock.Anything, "kc-fail").Return(errors.New("keycloak error"))

		svc := NewService(mockRepo, mockAuth)
		err := svc.CheckAndCleanInconsistentState(context.Background(), "fail@example.com")

		if err != domain.ErrKeycloakCleanupFailed {
			t.Fatalf("expected ErrKeycloakCleanupFailed, got %v", err)
		}
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_GetUserByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		kcUser := &gocloak.User{ID: strPtr("kc-123"), Email: strPtr("test@example.com")}
		mockAuth.On("GetUserByEmail", mock.Anything, "test@example.com").Return(kcUser, nil)

		svc := NewService(nil, mockAuth)
		result, err := svc.GetUserByEmail(context.Background(), "test@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if *result.ID != "kc-123" {
			t.Fatalf("expected ID kc-123, got %s", *result.ID)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("GetUserByEmail", mock.Anything, "unknown@example.com").Return(nil, errors.New("not found"))

		svc := NewService(nil, mockAuth)
		_, err := svc.GetUserByEmail(context.Background(), "unknown@example.com")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockAuth.AssertExpectations(t)
	})
}

func boolPtr(b bool) *bool { return &b }

func TestEmployeeService_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		kcUser := &gocloak.User{
			ID:            strPtr("kc-123"),
			Email:         strPtr("test@example.com"),
			EmailVerified: boolPtr(true),
		}
		token := &gocloak.JWT{AccessToken: "access-token"}

		mockAuth.On("GetUserByEmail", mock.Anything, "test@example.com").Return(kcUser, nil)
		mockAuth.On("LoginUser", mock.Anything, "test@example.com", "password123").Return(token, nil)

		svc := NewService(nil, mockAuth)
		result, err := svc.Login(context.Background(), "test@example.com", "password123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.AccessToken != "access-token" {
			t.Fatalf("expected access-token, got %s", result.AccessToken)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("GetUserByEmail", mock.Anything, "unknown@example.com").Return(nil, errors.New("not found"))

		svc := NewService(nil, mockAuth)
		_, err := svc.Login(context.Background(), "unknown@example.com", "password123")

		if err != domain.ErrUserNotFound {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("email not verified", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		kcUser := &gocloak.User{
			ID:            strPtr("kc-123"),
			Email:         strPtr("unverified@example.com"),
			EmailVerified: boolPtr(false),
		}
		mockAuth.On("GetUserByEmail", mock.Anything, "unverified@example.com").Return(kcUser, nil)
		mockAuth.On("SendVerificationEmail", mock.Anything, "kc-123").Return(nil)

		svc := NewService(nil, mockAuth)
		_, err := svc.Login(context.Background(), "unverified@example.com", "password123")

		if err != domain.ErrorEmailNotVerified {
			t.Fatalf("expected ErrorEmailNotVerified, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("login error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		kcUser := &gocloak.User{
			ID:            strPtr("kc-123"),
			Email:         strPtr("test@example.com"),
			EmailVerified: boolPtr(true),
		}
		mockAuth.On("GetUserByEmail", mock.Anything, "test@example.com").Return(kcUser, nil)
		mockAuth.On("LoginUser", mock.Anything, "test@example.com", "wrongpass").Return(nil, errors.New("invalid credentials"))

		svc := NewService(nil, mockAuth)
		_, err := svc.Login(context.Background(), "test@example.com", "wrongpass")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_SendVerificationEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SendVerificationEmail", mock.Anything, "kc-123").Return(nil)

		svc := NewService(nil, mockAuth)
		err := svc.SendVerificationEmail(context.Background(), "kc-123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SendVerificationEmail", mock.Anything, "kc-fail").Return(errors.New("email error"))

		svc := NewService(nil, mockAuth)
		err := svc.SendVerificationEmail(context.Background(), "kc-fail")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_SendPasswordResetEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SendPasswordResetEmail", mock.Anything, "test@example.com").Return(nil)

		svc := NewService(nil, mockAuth)
		err := svc.SendPasswordResetEmail(context.Background(), "test@example.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("SendPasswordResetEmail", mock.Anything, "fail@example.com").Return(errors.New("reset error"))

		svc := NewService(nil, mockAuth)
		err := svc.SendPasswordResetEmail(context.Background(), "fail@example.com")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_UpdatePassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("ValidateActionToken", mock.Anything, "valid-token").Return("kc-123", "test@example.com", nil)
		mockAuth.On("SetPassword", mock.Anything, "kc-123", "newpass", false).Return(nil)

		svc := NewService(nil, mockAuth)
		email, err := svc.UpdatePassword(context.Background(), "valid-token", "newpass")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if email != "test@example.com" {
			t.Fatalf("expected test@example.com, got %s", email)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("ValidateActionToken", mock.Anything, "bad-token").Return("", "", errors.New("invalid"))

		svc := NewService(nil, mockAuth)
		_, err := svc.UpdatePassword(context.Background(), "bad-token", "newpass")

		if err != domain.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("set password error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("ValidateActionToken", mock.Anything, "valid-token").Return("kc-123", "test@example.com", nil)
		mockAuth.On("SetPassword", mock.Anything, "kc-123", "newpass", false).Return(errors.New("password error"))

		svc := NewService(nil, mockAuth)
		_, err := svc.UpdatePassword(context.Background(), "valid-token", "newpass")

		if err != domain.ErrPasswordUpdateFailed {
			t.Fatalf("expected ErrPasswordUpdateFailed, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_ChangePassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		kcUser := &gocloak.User{ID: strPtr("kc-123"), Email: strPtr("test@example.com")}
		token := &gocloak.JWT{AccessToken: "token"}

		mockAuth.On("GetUserByEmail", mock.Anything, "test@example.com").Return(kcUser, nil)
		mockAuth.On("LoginUser", mock.Anything, "test@example.com", "oldpass").Return(token, nil)
		mockAuth.On("SetPassword", mock.Anything, "kc-123", "newpass", false).Return(nil)

		svc := NewService(nil, mockAuth)
		email, err := svc.ChangePassword(context.Background(), "test@example.com", "oldpass", "newpass")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if email != "test@example.com" {
			t.Fatalf("expected test@example.com, got %s", email)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		mockAuth.On("GetUserByEmail", mock.Anything, "unknown@example.com").Return(nil, errors.New("not found"))

		svc := NewService(nil, mockAuth)
		_, err := svc.ChangePassword(context.Background(), "unknown@example.com", "oldpass", "newpass")

		if err != domain.ErrUserNotFound {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("invalid current password", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		kcUser := &gocloak.User{ID: strPtr("kc-123"), Email: strPtr("test@example.com")}

		mockAuth.On("GetUserByEmail", mock.Anything, "test@example.com").Return(kcUser, nil)
		mockAuth.On("LoginUser", mock.Anything, "test@example.com", "wrongpass").Return(nil, errors.New("invalid creds"))

		svc := NewService(nil, mockAuth)
		_, err := svc.ChangePassword(context.Background(), "test@example.com", "wrongpass", "newpass")

		if err != domain.ErrInvalidCurrentPassword {
			t.Fatalf("expected ErrInvalidCurrentPassword, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("set password error", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		kcUser := &gocloak.User{ID: strPtr("kc-123"), Email: strPtr("test@example.com")}
		token := &gocloak.JWT{AccessToken: "token"}

		mockAuth.On("GetUserByEmail", mock.Anything, "test@example.com").Return(kcUser, nil)
		mockAuth.On("LoginUser", mock.Anything, "test@example.com", "oldpass").Return(token, nil)
		mockAuth.On("SetPassword", mock.Anything, "kc-123", "newpass", false).Return(errors.New("password error"))

		svc := NewService(nil, mockAuth)
		_, err := svc.ChangePassword(context.Background(), "test@example.com", "oldpass", "newpass")

		if err != domain.ErrPasswordUpdateFailed {
			t.Fatalf("expected ErrPasswordUpdateFailed, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})
}

func TestEmployeeService_VerifyEmailByToken(t *testing.T) {
	// Helper: build a fake JWT with 3 dots-separated base64 segments
	mkToken := func(email string) string {
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256"}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(`{"email":"` + email + `"}`))
		sig := base64.RawURLEncoding.EncodeToString([]byte("sig"))
		return header + "." + payload + "." + sig
	}

	t.Run("success", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		email := "test@example.com"
		notVerified := false
		kcUser := &gocloak.User{ID: strPtr("kc-123"), Email: strPtr(email), EmailVerified: &notVerified}

		mockAuth.On("GetUserByEmail", mock.Anything, email).Return(kcUser, nil)
		mockAuth.On("VerifyEmail", mock.Anything, "kc-123").Return(nil)

		svc := NewService(nil, mockAuth)
		result, err := svc.VerifyEmailByToken(context.Background(), mkToken(email))

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result != email {
			t.Fatalf("expected %s, got %s", email, result)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("invalid token format", func(t *testing.T) {
		svc := NewService(nil, nil)
		_, err := svc.VerifyEmailByToken(context.Background(), "not-a-jwt")

		if err != domain.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("user not found in keycloak", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		email := "unknown@example.com"

		mockAuth.On("GetUserByEmail", mock.Anything, email).Return(nil, errors.New("not found"))

		svc := NewService(nil, mockAuth)
		_, err := svc.VerifyEmailByToken(context.Background(), mkToken(email))

		if err != domain.ErrUserNotFound {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("email already verified", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		email := "verified@example.com"
		verified := true
		kcUser := &gocloak.User{ID: strPtr("kc-456"), Email: strPtr(email), EmailVerified: &verified}

		mockAuth.On("GetUserByEmail", mock.Anything, email).Return(kcUser, nil)

		svc := NewService(nil, mockAuth)
		result, err := svc.VerifyEmailByToken(context.Background(), mkToken(email))

		if err != domain.ErrEmailAlreadyVerified {
			t.Fatalf("expected ErrEmailAlreadyVerified, got %v", err)
		}
		if result != email {
			t.Fatalf("expected %s, got %s", email, result)
		}
		mockAuth.AssertExpectations(t)
	})

	t.Run("VerifyEmail fails", func(t *testing.T) {
		mockAuth := new(mocks.MockAuthClient)
		email := "test@example.com"
		notVerified := false
		kcUser := &gocloak.User{ID: strPtr("kc-789"), Email: strPtr(email), EmailVerified: &notVerified}

		mockAuth.On("GetUserByEmail", mock.Anything, email).Return(kcUser, nil)
		mockAuth.On("VerifyEmail", mock.Anything, "kc-789").Return(errors.New("kc error"))

		svc := NewService(nil, mockAuth)
		_, err := svc.VerifyEmailByToken(context.Background(), mkToken(email))

		if err == nil {
			t.Fatal("expected error")
		}
		mockAuth.AssertExpectations(t)
	})
}
