package tests

import (
	"log"
	"testing"

	"github.com/deejcoder/spidernet-api/storage"
	"github.com/deejcoder/spidernet-api/util/config"
)

func TestCreateServer(t *testing.T) {
	config.InitConfig()
	client := storage.Connect()

	// insert a new server and test for existance
	tags := []string{"Web", "Ubuntu"}
	exists, err := storage.CreateServer(client, "192.168.100.6:80", tags)
	if err != nil {
		log.Panic(err)
	}

	if exists {
		log.Panic("Doc already exists")
	}

	exists = storage.ServerExists(client, "192.168.100.6:80")
	if !exists {
		log.Panic("Doc doesn't exist when it should!")
	}

	// insert another and test that it doesn't overwrite an existing one
	exists, err = storage.CreateServer(client, "192.168.100.6:80", nil)
	if err != nil {
		log.Panic(err)
	}

	if !exists {
		log.Panic("Doc doesn't already exist when it should!")
	}

}
