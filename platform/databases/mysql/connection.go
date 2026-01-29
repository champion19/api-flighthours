package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/champion19/api-flighthours/config"
	 "github.com/champion19/api-flighthours/platform/logger"
	_ "github.com/go-sql-driver/mysql"
)

var log logger.Logger = logger.NewSlogLogger()


func GetDB(dbConfig config.Database) (*sql.DB, error) {
	log.Info(logger.LogDBConnecting,
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
		log.Debug(logger.LogDBSSLEnabled, "tsl", dbConfig.SSL)
		dsn += "&tls=" + dbConfig.SSL
	}

	db, err := sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		log.Error(logger.LogDBConnectionError,
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error to connect to database: %w", err)
	}

	log.Debug(logger.LogDBPoolConfig,
		"max_open_conns", dbConfig.MaxOpenConns,
		"max_idle_conns", dbConfig.MaxIdleConns,
		"conn_max_lifetime", dbConfig.ConnMaxLifetime,
		"conn_max_idle_time", dbConfig.ConnMaxIdleTime,
	)

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime))
	db.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime))

	log.Info(logger.LogDBPinging)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Error(logger.LogDBPingError,
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	log.Success(logger.LogDBConnected,
		"host", dbConfig.Host,
		"database", dbConfig.Name)

	return db, nil
}
