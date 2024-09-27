package database

import (
	"database/sql"
	"embed"
	"fmt"

	"project-template/config"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var EmbeddedMigrations embed.FS

func init() {
	pgx.ErrNoRows = sql.ErrNoRows
}

func New(c *config.Database) (*sql.DB, error) {
	defaultDSN := fmt.Sprintf("postgresql://%v:%v@%v:%v/postgres?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.SSLMode)
	defaultDB, err := sql.Open("pgx", defaultDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to default database: %w", err)
	}
	defer defaultDB.Close()

	// 1. Get existing databases
	var exists bool
	err = defaultDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", c.Database).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	// 2. Create the database if it does not exist
	if !exists {
		_, err = defaultDB.Exec("CREATE DATABASE " + c.Database)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	// 3. Create the connection pool
	dsn := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to specified database: %w", err)
	}

	// 4. Ping the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 5. Run the migrations
	goose.SetBaseFS(EmbeddedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	// 6. Return the database
	return db, nil
}
