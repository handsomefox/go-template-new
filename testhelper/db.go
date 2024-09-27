package testhelper

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"project-template/config"
	"project-template/database"

	"github.com/ory/dockertest/v3"
)

type TestDatabase struct {
	DB       *sql.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func SetupTestDatabase() (*TestDatabase, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	resource, err := pool.Run("postgres", "13", []string{"POSTGRES_PASSWORD=secret"})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %w", err)
	}

	var db *sql.DB
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("pgx", fmt.Sprintf("postgresql://postgres:secret@localhost:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	portStr := resource.GetPort("5432/tcp")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse port number: %w", err)
	}

	cfg := &config.Database{
		Host:     "localhost",
		Port:     port,
		Username: "postgres",
		Password: "secret",
		Database: "postgres",
		SSLMode:  "disable",
	}

	db, err = database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not set up database: %w", err)
	}

	return &TestDatabase{
		DB:       db,
		pool:     pool,
		resource: resource,
	}, nil
}

func (td *TestDatabase) TearDown() {
	if err := td.pool.Purge(td.resource); err != nil {
		log.Printf("Could not purge resource: %s", err)
	}
}
