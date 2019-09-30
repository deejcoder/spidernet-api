package client

import (
	"database/sql"
	"fmt"

	config "github.com/deejcoder/spidernet-api/util/config"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	sb "github.com/huandu/go-sqlbuilder"
)

type PostgresInstance struct {
	Db *sql.DB
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

	instance.Db = db
	return nil
}

// Migrate checks if there are any postgres changes, updates if so
func (instance *PostgresInstance) Migrate() error {
	driver, err := postgres.WithInstance(instance.Db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://storage/client/migrations", "postgres", driver)
	if err != nil {
		return err
	}
	m.Steps(2)
	return nil
}

// Update updates an existing tuple; this is a wrapper for go-sqlbuilder
func (instance *PostgresInstance) Update(table string, s *sb.Struct, value interface{}) error {
	ub := s.Update(table, value)
	sql, args := ub.BuildWithFlavor(sb.PostgreSQL)
	_, err := instance.Db.Query(sql, args...)
	return err
}

// Delete deletes an existing tuple WHERE key is value; this is a wrapper for go-sqlbuilder
func (instance *PostgresInstance) Delete(table string, key string, value interface{}) error {
	db := sb.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom(table)
	db.Where(db.Equal(key, value))
	sql, args := db.Build()

	_, err := instance.Db.Query(sql, args...)
	return err
}

func (instance *PostgresInstance) GetOne(table string, key string, value interface{}, dest interface{}) error {
	selectb := sb.PostgreSQL.NewSelectBuilder()
	selectb.Select("*")
	selectb.From(table)
	selectb.Where(selectb.E(key, value))
	selectb.Limit(1)

	sql, args := selectb.Build()
	res := instance.Db.QueryRow(sql, args...)

	err := res.Scan(&dest)
	if err != nil {
		return err
	}
	return nil
}
