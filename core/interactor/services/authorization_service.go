package services

/*¨
import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports"


)

type authorizationService struct {
	keycloakClient ports.AuthClient
	repository     ports.Repository
}


func NewAuthorizationService(keycloakClient ports.AuthClient, repository ports.Repository) ports.AuthorizationService {
	return &authorizationService{
		keycloakClient: keycloakClient,
		repository:     repository,
	}
}

func (a *authorizationService) SyncUserToKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {

	if employee.KeycloakUserID != "" {
		return employee.KeycloakUserID, nil
	}


	keycloakUserID, err := a.keycloakClient.CreateUser(ctx, employee)
	if err != nil {

		existingUser, getErr := a.keycloakClient.GetUserByEmail(ctx, employee.Email)
		if getErr != nil {
			return "", fmt.Errorf("failed to create or get user in keycloak: %w", err)
		}
		keycloakUserID = *existingUser.ID
	}

	employee.KeycloakUserID = keycloakUserID
	err = a.repository.UpdateEmployee(*employee)
	if err != nil {

		_ = a.keycloakClient.DeleteUser(ctx, keycloakUserID)
		return "", fmt.Errorf("failed to update person with keycloak user id: %w", err)
	}

	return keycloakUserID, nil
}


func (a *authorizationService) DeleteUserFromKeycloak(ctx context.Context, keycloakUserID string) error {
	if keycloakUserID == "" {
		return fmt.Errorf("keycloakUserID cannot be empty")
	}

	err := a.keycloakClient.DeleteUser(ctx, keycloakUserID)
	if err != nil {
		return fmt.Errorf("failed to delete user from keycloak: %w", err)
	}

	return nil
}


func (a *authorizationService) LoginUser(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password cannot be empty")
	}

	token, err := a.keycloakClient.LoginUser(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return token, nil
}


func (a *authorizationService) SetUserPassword(ctx context.Context, keycloakUserID string, password string) error {
	if keycloakUserID == "" || password == "" {
		return fmt.Errorf("keycloakUserID and password cannot be empty")
	}


	err := a.keycloakClient.SetPassword(ctx, keycloakUserID, password, false)
	if err != nil {
		return fmt.Errorf("failed to set password in keycloak: %w", err)
	}

	return nil
}


func (a *authorizationService) AssignRole(ctx context.Context, employeeID string, roleName string) error {

	employee, err := a.repository.GetEmployeeByID(employeeID)
	if err != nil {
		return fmt.Errorf("employee not found: %w", err)
	}


	keycloakUserID, err := a.SyncUserToKeycloak(ctx, employee)
	if err != nil {
		return fmt.Errorf("failed to sync user to keycloak: %w", err)
	}


	err = a.keycloakClient.AssignRole(ctx, keycloakUserID, roleName)
	if err != nil {
		return fmt.Errorf("failed to assign role in keycloak: %w", err)
	}


	employee.Role = roleName
	err = a.repository.UpdateEmployee(*employee)
	if err != nil {
		fmt.Printf("Warning: failed to update role in local database: %v\n", err)
	}
	return nil
}


func (a *authorizationService) RemoveRole(ctx context.Context, employeeID string, roleName string) error {
	keycloakUserID, err := a.getKeycloakUserID(ctx, employeeID)
	if err != nil {
		return err
	}

	return a.keycloakClient.RemoveRole(ctx, keycloakUserID, roleName)
}


func (a *authorizationService) GetUserRoles(ctx context.Context, employeeID string) ([]string, error) {
	keycloakUserID, err := a.getKeycloakUserID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	roles, err := a.keycloakClient.GetUserRoles(ctx, keycloakUserID)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.Name != nil {
			roleNames = append(roleNames, *role.Name)
		}
	}

	return roleNames, nil
}


func (a *authorizationService) HasRole(ctx context.Context, employeeID string, roleName string) (bool, error) {
	roles, err := a.GetUserRoles(ctx, employeeID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role == roleName {
			return true, nil
		}
	}

	return false, nil
}


func (a *authorizationService) HasPermission(ctx context.Context, employeeID string, resource, action string) (bool, error) {

	roles, err := a.GetUserRoles(ctx, employeeID)
	if err != nil {
		return false, err
	}


	for _, role := range roles {
		switch role {
		case "admin":
			return true, nil
		case "moderator":
			if resource == "users" && (action == "read" || action == "update") {
				return true, nil
			}
		case "user":
			if resource == "profile" && action == "read" {
				return true, nil
			}
		}
	}

	return false, nil
}


func (a *authorizationService) CreateRole(ctx context.Context, roleName, description string) error {
	roles, err := a.keycloakClient.GetUserRoles(ctx, "dummy")
	if err == nil {
		for _, role := range roles {
			if role.Name != nil && *role.Name == roleName {
				return fmt.Errorf("role %s already exists", roleName)
			}
		}
	}
	return fmt.Errorf("create role functionality needs to be implemented in KeycloakClient")
}


func (a *authorizationService) GetAllRoles(ctx context.Context) ([]*string, error) {
	return nil, fmt.Errorf("get all roles functionality needs to be implemented in KeycloakClient")
}

func (a *authorizationService) getKeycloakUserID(ctx context.Context, employeeID string) (string, error) {
	employee, err := a.repository.GetEmployeeByID(employeeID)
	if err != nil {
		return "", fmt.Errorf("employee not found: %w", err)
	}


	if employee.KeycloakUserID != "" {
		return employee.KeycloakUserID, nil
	}
	return a.SyncUserToKeycloak(ctx, employee)
}
*/
