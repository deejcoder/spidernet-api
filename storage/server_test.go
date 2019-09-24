package storage

import (
	"log"
	"testing"

	"github.com/deejcoder/spidernet-api/util/config"
)

func setup() (*config.Config, *PostgresInstance) {
	cfg := config.InitConfig("../")

	instance := NewPostgresInstance()
	if err := instance.Connect(); err != nil {
		log.Fatal(err)
	}

	return cfg, instance
}

func TestServerOperations(t *testing.T) {
	_, instance := setup()
	tags := []string{"Router", "Gateway"}
	err := instance.CreateServer("192.168.20.1", "Modem Gateway", tags)
	if err != nil {
		log.Fatal(err)
	}
}
