package services

import (
	"context"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/output"
)

type CrewMemberService struct {
	repo output.CrewMemberRepository
}

func NewCrewMemberService(repo output.CrewMemberRepository) *CrewMemberService {
	return &CrewMemberService{repo: repo}
}

func (s *CrewMemberService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

func (s *CrewMemberService) SearchCrewMembers(ctx context.Context, employeeID, query string) ([]domain.CrewMember, error) {
	return s.repo.SearchCrewMembers(ctx, employeeID, query)
}

func (s *CrewMemberService) FindOrCreateCrewMemberTx(ctx context.Context, tx output.Tx, employeeID, name string, bp *string) (*domain.CrewMember, error) {
	return s.repo.FindOrCreateCrewMember(ctx, tx, employeeID, name, bp)
}

func (s *CrewMemberService) GetCrewMemberByID(ctx context.Context, id string) (*domain.CrewMember, error) {
	return s.repo.GetCrewMemberByID(ctx, id)
}
