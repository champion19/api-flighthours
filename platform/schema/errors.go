package schema

import (
	

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)



var (
	ErrBadRequest            = domain.ErrSchemaBadRequest
	ErrInvalidRequest        = domain.ErrSchemaInvalidRequest
	ErrSchemaReadFailed      = domain.ErrSchemaReadFailed
	ErrSchemaEmpty           = domain.ErrSchemaEmpty
	ErrSchemaCompileFailed   = domain.ErrSchemaCompileFailed
	ErrValidationFailed      = domain.ErrSchemaValidationFailed
	ErrBodyReadFailed        = domain.ErrSchemaBodyReadFailed
	ErrFieldPropertyMismatch = domain.ErrSchemaFieldFormat
	ErrFieldRequired         = domain.ErrSchemaFieldRequired
	ErrFieldTypeInvalid      = domain.ErrSchemaFieldType
	ErrMultipleFields        = domain.ErrSchemaMultipleFields
)
