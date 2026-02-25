package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/marcboeker/go-duckdb" // DuckDB driver
	_ "github.com/mattn/go-sqlite3"     // SQLite driver
)

// DBType represents the type of database
type DBType string

const (
	SQLite  DBType = "sqlite"
	DuckDB  DBType = "duckdb"
)

// DBConfig holds database configuration
type DBConfig struct {
	Type     DBType
	FilePath string
}

// Database interface for common operations
type Database interface {
	GetDB() *sql.DB
	Setup() error
	Type() DBType
	Close() error
}

// database implements the Database interface
type database struct {
	db       *sql.DB
	dbType   DBType
	filePath string
}

// NewDatabase creates a new database connection
func NewDatabase(config DBConfig) (Database, error) {
	var db *sql.DB
	var err error

	switch config.Type {
	case SQLite:
		db, err = sql.Open("sqlite3", config.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite: %w", err)
		}
	case DuckDB:
		db, err = sql.Open("duckdb", config.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open DuckDB: %w", err)
		}
		// Install inet extension for IP operations
		if _, err := db.Exec("INSTALL inet; LOAD inet;"); err != nil {
			log.Printf("Warning: Failed to load inet extension: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	return &database{
		db:       db,
		dbType:   config.Type,
		filePath: config.FilePath,
	}, nil
}

func (d *database) GetDB() *sql.DB {
	return d.db
}

func (d *database) Type() DBType {
	return d.dbType
}

func (d *database) Close() error {
	return d.db.Close()
}

func (d *database) Setup() error {
	var createSQL string

	switch d.dbType {
	case SQLite:
		createSQL = `
		CREATE TABLE IF NOT EXISTS net (
			net TEXT PRIMARY KEY,
			start_ip TEXT NOT NULL,
			end_ip TEXT NOT NULL,
			url TEXT NOT NULL,
			cloudplatform TEXT NOT NULL,
			iptype TEXT NOT NULL
		);`
	case DuckDB:
		createSQL = `
		CREATE TABLE IF NOT EXISTS net (
			net VARCHAR PRIMARY KEY,
			start_ip VARCHAR NOT NULL,
			end_ip VARCHAR NOT NULL,
			url VARCHAR NOT NULL,
			cloudplatform VARCHAR NOT NULL,
			iptype VARCHAR NOT NULL
		);`
	default:
		return fmt.Errorf("unsupported database type: %s", d.dbType)
	}

	_, err := d.db.Exec(createSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
