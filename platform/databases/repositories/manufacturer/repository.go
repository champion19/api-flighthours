package manufacturer

import(
	"database/sql"

	"github.com/champion19/api-flighthours/platform/logger"
)

const (
	QueryByID   = "SELECT id, name FROM manufacturer WHERE id = ? LIMIT 1"
	QueryGetAll = "SELECT id, name FROM manufacturer ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID *sql.Stmt
	stmtGetAll  *sql.Stmt
	db          *sql.DB
}


func NewManufacturerRepository(db *sql.DB) (*repository, error) {
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

	return &repository{
		db:          db,
		stmtGetByID: stmtGetByID,
		stmtGetAll:  stmtGetAll,
	}, nil
}


