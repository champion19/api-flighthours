package interactor

import(
	"context"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	)



	// ManufacturerInteractor orchestrates manufacturer operations
type ManufacturerInteractor struct {
	service input.ManufacturerService
}

// NewManufacturerInteractor creates a new manufacturer interactor
func NewManufacturerInteractor(service input.ManufacturerService) *ManufacturerInteractor {
	return &ManufacturerInteractor{
		service: service,
		}
}

// GetManufacturerByID retrieves a manufacturer by its ID
func (i *ManufacturerInteractor) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogManufacturerGet, "manufacturer_id", id)

	manufacturer, err := i.service.GetManufacturerByID(ctx, id)
	if err != nil {
		log.Error(logger.LogManufacturerGetError, "manufacturer_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogManufacturerGetOK, manufacturer.ToLogger())
	return manufacturer, nil
}

// ListManufacturers retrieves all manufacturers
func (i *ManufacturerInteractor) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogManufacturerList)

	manufacturers, err := i.service.ListManufacturers(ctx)
	if err != nil {
		log.Error(logger.LogManufacturerListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogManufacturerListOK, "count", len(manufacturers))
	return manufacturers, nil
}
