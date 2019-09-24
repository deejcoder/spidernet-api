package storage

import (
	"database/sql"
	"fmt"

	config "github.com/deejcoder/spidernet-api/util/config"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresInstance struct {
	client *sql.DB
}

func NewPostgresInstance() *PostgresInstance {
	return &PostgresInstance{}
}

// Connect connects to a postgres database
func (instance *PostgresInstance) Connect() error {
	// build connection string
	conf := config.GetConfig()
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Db,
		conf.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	// test connection
	if err = db.Ping(); err != nil {
		return err
	}

	instance.client = db
	return nil
}

// Migrate checks if there are any postgres changes, updates if so
func (instance *PostgresInstance) Migrate() error {
	driver, err := postgres.WithInstance(instance.client, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://storage/migrations", "postgres", driver)
	if err != nil {
		return err
	}
	m.Steps(2)
	return nil
}
