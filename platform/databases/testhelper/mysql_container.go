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

// SetupMessageSchema creates the system_messages table for testing
func (m *MySQLContainer) SetupMessageSchema(ctx context.Context) error {
	schema := `
		CREATE TABLE IF NOT EXISTS system_messages (
			id VARCHAR(36) NOT NULL,
			message_code VARCHAR(50) NOT NULL,
			type VARCHAR(20) NOT NULL,
			category VARCHAR(50),
			module VARCHAR(50),
			message_title VARCHAR(100),
			message_content TEXT,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY(id)
		);
	`
	_, err := m.DB.ExecContext(ctx, schema)
	return err
}

// InsertMessage inserts test data into system_messages table
func (m *MySQLContainer) InsertMessage(ctx context.Context, id, code, msgType, category, module, title, content string, active bool) error {
	_, err := m.DB.ExecContext(ctx,
		"INSERT INTO system_messages (id, message_code, type, category, module, message_title, message_content, is_active) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		id, code, msgType, category, module, title, content, active)
	return err
}

// CleanMessageTable removes all data from system_messages table
func (m *MySQLContainer) CleanMessageTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM system_messages")
	return err
}

// SetupEmployeeSchema creates the employee table for testing (requires airline table)
func (m *MySQLContainer) SetupEmployeeSchema(ctx context.Context) error {
	// First ensure airline table exists
	if err := m.SetupAirlineSchema(ctx); err != nil {
		return err
	}
	schema := `
		CREATE TABLE IF NOT EXISTS employee (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(50) NOT NULL,
			airline VARCHAR(36) NULL,
			email VARCHAR(150) NOT NULL UNIQUE,
			identification_number VARCHAR(10) NOT NULL,
			bp VARCHAR(16) NULL,
			start_date DATE NOT NULL,
			end_date DATE,
			active BOOLEAN NOT NULL,
			role VARCHAR(10) NOT NULL,
			keycloak_user_id VARCHAR(36) NULL,
			PRIMARY KEY(id)
		);
	`
	_, err := m.DB.ExecContext(ctx, schema)
	return err
}

// InsertEmployee inserts test data into employee table
func (m *MySQLContainer) InsertEmployee(ctx context.Context, id, name, email, identNum, role string, active bool) error {
	_, err := m.DB.ExecContext(ctx,
		"INSERT INTO employee (id, name, email, identification_number, start_date, active, role) VALUES (?, ?, ?, ?, CURDATE(), ?, ?)",
		id, name, email, identNum, active, role)
	return err
}

// CleanEmployeeTable removes all data from employee table
func (m *MySQLContainer) CleanEmployeeTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM employee")
	return err
}

// ===========================================
// Airport Schema (base table for routes)
// ===========================================

// SetupAirportSchema creates the airport table for testing
func (m *MySQLContainer) SetupAirportSchema(ctx context.Context) error {
	schema := `
		CREATE TABLE IF NOT EXISTS airport (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(50) NOT NULL,
			iata_code VARCHAR(3),
			status BOOLEAN NOT NULL,
			airport_type VARCHAR(13),
			PRIMARY KEY(id)
		);
	`
	_, err := m.DB.ExecContext(ctx, schema)
	return err
}

// InsertAirport inserts test data into airport table
func (m *MySQLContainer) InsertAirport(ctx context.Context, id, name, iataCode, airportType string, status bool) error {
	_, err := m.DB.ExecContext(ctx,
		"INSERT INTO airport (id, name, iata_code, status, airport_type) VALUES (?, ?, ?, ?, ?)",
		id, name, iataCode, status, airportType)
	return err
}

// CleanAirportTable removes all data from airport table
func (m *MySQLContainer) CleanAirportTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM airport")
	return err
}

// ===========================================
// Route Schema (depends on airport)
// ===========================================

// SetupRouteSchema creates the route table for testing (requires airport table)
func (m *MySQLContainer) SetupRouteSchema(ctx context.Context) error {
	// First ensure airport table exists
	if err := m.SetupAirportSchema(ctx); err != nil {
		return err
	}
	schema := `
		CREATE TABLE IF NOT EXISTS route (
			id VARCHAR(36) NOT NULL,
			origin_airport_id VARCHAR(36) NOT NULL,
			destination_airport_id VARCHAR(36) NOT NULL,
			airport_type VARCHAR(13) NOT NULL,
			estimated_flight_time TIME,
			PRIMARY KEY (id),
			CONSTRAINT fk_route_origin_airport FOREIGN KEY (origin_airport_id) REFERENCES airport(id),
			CONSTRAINT fk_route_destination_airport FOREIGN KEY (destination_airport_id) REFERENCES airport(id)
		);
	`
	_, err := m.DB.ExecContext(ctx, schema)
	return err
}

// InsertRoute inserts test data into route table
func (m *MySQLContainer) InsertRoute(ctx context.Context, id, originAirportID, destAirportID, airportType, estimatedTime string) error {
	_, err := m.DB.ExecContext(ctx,
		"INSERT INTO route (id, origin_airport_id, destination_airport_id, airport_type, estimated_flight_time) VALUES (?, ?, ?, ?, ?)",
		id, originAirportID, destAirportID, airportType, estimatedTime)
	return err
}

// CleanRouteTable removes all data from route table
func (m *MySQLContainer) CleanRouteTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM route")
	return err
}

// ===========================================
// Airline Route Schema (depends on route and airline)
// ===========================================

// SetupAirlineRouteSchema creates the airline_route table for testing (requires route and airline tables)
func (m *MySQLContainer) SetupAirlineRouteSchema(ctx context.Context) error {
	// First ensure dependent tables exist
	if err := m.SetupAirlineSchema(ctx); err != nil {
		return err
	}
	if err := m.SetupRouteSchema(ctx); err != nil {
		return err
	}
	schema := `
		CREATE TABLE IF NOT EXISTS airline_route (
			id VARCHAR(36) NOT NULL,
			route_id VARCHAR(36) NOT NULL,
			airline_id VARCHAR(36) NOT NULL,
			status BOOLEAN NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_airline_route_route FOREIGN KEY (route_id) REFERENCES route(id),
			CONSTRAINT fk_airline_route_airline FOREIGN KEY (airline_id) REFERENCES airline(id)
		);
	`
	_, err := m.DB.ExecContext(ctx, schema)
	return err
}

// InsertAirlineRoute inserts test data into airline_route table
func (m *MySQLContainer) InsertAirlineRoute(ctx context.Context, id, routeID, airlineID string, status bool) error {
	_, err := m.DB.ExecContext(ctx,
		"INSERT INTO airline_route (id, route_id, airline_id, status) VALUES (?, ?, ?, ?)",
		id, routeID, airlineID, status)
	return err
}

// CleanAirlineRouteTable removes all data from airline_route table
func (m *MySQLContainer) CleanAirlineRouteTable(ctx context.Context) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM airline_route")
	return err
}
