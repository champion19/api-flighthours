package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/champion19/flighthours-api/config"
	loggerPkg "github.com/champion19/flighthours-api/platform/logger"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB(dbConfig config.Database, logger loggerPkg.Logger) (*sql.DB, error) {
	logger.Info(loggerPkg.LogDBConnecting,
		"host", dbConfig.Host,
		"port", dbConfig.Port,
		"database", dbConfig.Name,
		"driver", dbConfig.Driver)

	var dsn string

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	if dbConfig.SSL != "" {
		logger.Debug(loggerPkg.LogDBSSLEnabled, "tsl", dbConfig.SSL)
		dsn += "&tls=" + dbConfig.SSL
	}

	db, err := sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		logger.Error(loggerPkg.LogDBConnectionError,
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error to connect to database: %w", err)
	}

	logger.Debug(loggerPkg.LogDBPoolConfig,
		"max_open_conns", dbConfig.MaxOpenConns,
		"max_idle_conns", dbConfig.MaxIdleConns,
		"conn_max_lifetime", dbConfig.ConnMaxLifetime,
		"conn_max_idle_time", dbConfig.ConnMaxIdleTime,
	)

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime))
	db.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime))

	logger.Info(loggerPkg.LogDBPinging)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		logger.Error(loggerPkg.LogDBPingError,
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	logger.Success(loggerPkg.LogDBConnected,
		"host", dbConfig.Host,
		"database", dbConfig.Name)

	return db, nil
}
