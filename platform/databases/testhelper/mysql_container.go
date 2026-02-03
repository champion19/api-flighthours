//go:build integration
// +build integration

package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MySQLContainer wraps a testcontainers MySQL instance
type MySQLContainer struct {
	Container testcontainers.Container
	DB        *sql.DB
	Host      string
	Port      string
}

// StartMySQL creates and starts a MySQL container for testing
func StartMySQL(ctx context.Context) (*MySQLContainer, error) {
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("flighthours_test"),
		mysql.WithUsername("test"),
		mysql.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("ready for connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start mysql container: %w", err)
	}

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}

	dsn := fmt.Sprintf("test:test@tcp(%s:%s)/flighthours_test?parseTime=true", host, port.Port())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	// Wait for connection
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	return &MySQLContainer{
		Container: mysqlContainer,
		DB:        db,
		Host:      host,
		Port:      port.Port(),
	}, nil
}

// Stop terminates the MySQL container
func (m *MySQLContainer) Stop(ctx context.Context) error {
	if m.DB != nil {
		m.DB.Close()
	}
	if m.Container != nil {
		return m.Container.Terminate(ctx)
	}
	return nil
}

// SetupAirlineSchema creates the airline table for testing
func (m *MySQLContainer) SetupAirlineSchema(ctx context.Context) error {
	schema := `
		CREATE TABLE IF NOT EXISTS airline (
			id VARCHAR(36) NOT NULL,
			airline_name VARCHAR(30) NOT NULL,
			airline_code VARCHAR(3) NOT NULL,
			status BOOLEAN NOT NULL,
			PRIMARY KEY (id)
		);
	`
	_, err := m.DB.ExecContext(ctx, schema)
	return err
}

// InsertAirline inserts test data into airline table
func (m *MySQLContainer) InsertAirline(ctx context.Context, id, name, code string, status bool) error {
	_, err := m.DB.ExecContext(ctx,
		"INSERT INTO airline (id, airline_name, airline_code, status) VALUES (?, ?, ?, ?)",
		id, name, code, status)
	return err
}

// CleanAirlineTable removes all data from airline table
func (m *MySQLContainer) CleanAirlineTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM airline")
	return err
}

// SetupEngineSchema creates the engine table for testing
func (m *MySQLContainer) SetupEngineSchema(ctx context.Context) error {
	schema := `
		CREATE TABLE IF NOT EXISTS engine (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(3) NOT NULL,
			PRIMARY KEY(id)
		);
	`
	_, err := m.DB.ExecContext(ctx, schema)
	return err
}

// InsertEngine inserts test data into engine table
func (m *MySQLContainer) InsertEngine(ctx context.Context, id, name string) error {
	_, err := m.DB.ExecContext(ctx,
		"INSERT INTO engine (id, name) VALUES (?, ?)",
		id, name)
	return err
}

// CleanEngineTable removes all data from engine table
func (m *MySQLContainer) CleanEngineTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM engine")
	return err
}
