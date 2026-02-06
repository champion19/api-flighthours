package aircraftmodel

import(
	
	"database/sql"

	"github.com/champion19/api-flighthours/platform/logger"
)


const (
	QueryByID            = "SELECT am.id, am.model_name, am.aircraft_type_name, e.name AS engine_type_name, am.family, m.name AS manufacturer, am.status FROM aircraft_model am LEFT JOIN engine e ON am.engine_type_id = e.id LEFT JOIN manufacturer m ON am.manufacturer_id = m.id WHERE am.id = ? LIMIT 1"
	QueryGetAll          = "SELECT am.id, am.model_name, am.aircraft_type_name, e.name AS engine_type_name, am.family, m.name AS manufacturer, am.status FROM aircraft_model am LEFT JOIN engine e ON am.engine_type_id = e.id LEFT JOIN manufacturer m ON am.manufacturer_id = m.id ORDER BY am.model_name"
	QueryGetByEngineType = "SELECT am.id, am.model_name, am.aircraft_type_name, e.name AS engine_type_name, am.family, m.name AS manufacturer, am.status FROM aircraft_model am LEFT JOIN engine e ON am.engine_type_id = e.id LEFT JOIN manufacturer m ON am.manufacturer_id = m.id WHERE e.name = ? ORDER BY am.model_name"
	QueryGetByFamily     = "SELECT am.id, am.model_name, am.aircraft_type_name, e.name AS engine_type_name, am.family, m.name AS manufacturer, am.status FROM aircraft_model am LEFT JOIN engine e ON am.engine_type_id = e.id LEFT JOIN manufacturer m ON am.manufacturer_id = m.id WHERE am.family = ? ORDER BY am.model_name"
	QueryUpdateStatus    = "UPDATE aircraft_model SET status = ? WHERE id = ?"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID         *sql.Stmt
	stmtGetAll          *sql.Stmt
	stmtGetByEngineType *sql.Stmt
	stmtGetByFamily     *sql.Stmt
	stmtUpdateStatus    *sql.Stmt
	db                  *sql.DB
}

// NewAircraftModelRepository creates a new aircraft model repository with prepared statements
func NewAircraftModelRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByEngineType, err := db.Prepare(QueryGetByEngineType)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByFamily, err := db.Prepare(QueryGetByFamily)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	return &repository{
		db:                  db,
		stmtGetByID:         stmtGetByID,
		stmtGetAll:          stmtGetAll,
		stmtGetByEngineType: stmtGetByEngineType,
		stmtGetByFamily:     stmtGetByFamily,
		stmtUpdateStatus:    stmtUpdateStatus,
	}, nil
}

