package handlers

import (
	"errors"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
)

func TestFromDomainEngine(t *testing.T) {
	t.Run("converts domain engine to response", func(t *testing.T) {
		engine := &domain.Engine{
			ID:   "raw-engine-id",
			Name: "CFM56-7B",
		}

		result := FromDomainEngine(engine, "encoded-engine-id")

		if result.ID != "encoded-engine-id" {
			t.Errorf("expected 'encoded-engine-id', got %q", result.ID)
		}
		if result.Name != "CFM56-7B" {
			t.Errorf("expected 'CFM56-7B', got %q", result.Name)
		}
	})
}

func TestToEngineListResponse(t *testing.T) {
	t.Run("converts engine slice to list response with links", func(t *testing.T) {
		engines := []domain.Engine{
			{ID: "engine-1", Name: "CFM56-7B"},
			{ID: "engine-2", Name: "LEAP-1A"},
		}

		encodeFunc := func(id string) (string, error) {
			return "enc-" + id, nil
		}

		result := ToEngineListResponse(engines, encodeFunc, "http://api.test")

		if result.Total != 2 {
			t.Errorf("expected total 2, got %d", result.Total)
		}
		if len(result.Engines) != 2 {
			t.Errorf("expected 2 engines, got %d", len(result.Engines))
		}
		if result.Engines[0].ID != "enc-engine-1" {
			t.Errorf("expected 'enc-engine-1', got %q", result.Engines[0].ID)
		}
		if result.Engines[0].Name != "CFM56-7B" {
			t.Errorf("expected 'CFM56-7B', got %q", result.Engines[0].Name)
		}
		// Check that HATEOAS links are added
		if len(result.Engines[0].Links) == 0 {
			t.Error("expected HATEOAS links on engine response")
		}
		if len(result.Links) == 0 {
			t.Error("expected collection-level HATEOAS links")
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		engines := []domain.Engine{}
		encodeFunc := func(id string) (string, error) {
			return "enc-" + id, nil
		}

		result := ToEngineListResponse(engines, encodeFunc, "http://api.test")

		if result.Total != 0 {
			t.Errorf("expected total 0, got %d", result.Total)
		}
	})

	t.Run("uses original ID when encoding fails", func(t *testing.T) {
		engines := []domain.Engine{
			{ID: "fallback-id", Name: "GE90"},
		}

		encodeFunc := func(id string) (string, error) {
			return "", errors.New("encoding failed")
		}

		result := ToEngineListResponse(engines, encodeFunc, "http://api.test")

		if result.Engines[0].ID != "fallback-id" {
			t.Errorf("expected 'fallback-id', got %q", result.Engines[0].ID)
		}
	})

	t.Run("works without base URL", func(t *testing.T) {
		engines := []domain.Engine{
			{ID: "engine-1", Name: "CFM56"},
		}

		encodeFunc := func(id string) (string, error) {
			return "enc-" + id, nil
		}

		result := ToEngineListResponse(engines, encodeFunc, "")

		if result.Engines[0].Links != nil {
			t.Error("expected no links with empty baseURL")
		}
		if result.Links != nil {
			t.Error("expected no collection links with empty baseURL")
		}
	})
}
