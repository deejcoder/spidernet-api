package storage

import (
	"fmt"

	"github.com/sirupsen/logrus"

	config "github.com/deejcoder/spidernet-api/util/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresInstance struct {
	Db *gorm.DB
}

func NewPostgresInstance() *PostgresInstance {
	return &PostgresInstance{}
}

// Connect connects to a postgres database
func (instance *PostgresInstance) Connect() error {
	// build connection string
	conf := config.GetConfig()

	connStr := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s",
		conf.Database.Host,
		conf.Database.User,
		conf.Database.Db,
		conf.Database.Password,
		conf.Database.SSLMode,
	)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// set the db
	instance.Db = db

	// setup extensions, auto migrations etc.
	instance.OnConnect()
	return nil
}

func (instance *PostgresInstance) OnConnect() {
	instance.Db.AutoMigrate(&Server{}, &Tag{}, &User{})

	// try add pg_trgm extension if it doesn't exist
	out := instance.Db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;")
	if err := out.Error; err != nil {
		logrus.Fatalf("Cannot create pg_trgm extension, %s", err)
	}
}
