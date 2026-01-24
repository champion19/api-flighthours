package dto

import (
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type RegisterEmployee struct {
	Employee domain.Employee `json:"employee"`
	Message  string          `json:"message"`
}

type UserSyncStatus struct {
	EmployeeID     string `json:"employee_id"`
	KeycloakUserID string `json:"keycloak_user_id"`
	IsSynced       bool   `json:"is_synced"`
	LastSyncAt     string `json:"last_sync_at,omitempty"`
}

type TokenResponse struct {
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}


func FromDomainToDTO(employee *domain.Employee) *RegisterEmployee {
	return &RegisterEmployee{
		Employee: *employee,
		Message:  "Employee located successfully",
	}
}

