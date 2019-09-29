package server

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/deejcoder/spidernet-api/storage/client"
	"github.com/deejcoder/spidernet-api/util/config"
	log "github.com/sirupsen/logrus"
)

func setup() (*config.Config, *client.PostgresInstance) {
	cfg := config.InitConfig("../../")

	// configure & connect to postgres
	psql := client.NewPostgresInstance()
	if err := psql.Connect(); err != nil {
		log.Fatal(err)
	}

	return cfg, psql
}

func TestServerOperations(t *testing.T) {
	_, psql := setup()

	// we want to manage the servers with the postgres client
	mgr := NewServerManager(psql)
	CreateServers(mgr)

}

// CreateServers will do some magic!
func CreateServers(mgr *ServerManager) {
	tags := []string{
		"minecraft", "runescape", "web", "dns", "proxy", "database", "rust",
		"nitrox", "osrs", "social", "email", "community", "forum", "storage",
	}

	src := rand.NewSource(time.Now().Unix())
	rnd := rand.New(src)

	for i := 0; i < 50; i++ {
		// generate an ip
		addr := fmt.Sprintf("135.158.3.%d", i)

		// randomly select two tags
		index := rnd.Intn(len(tags)-1) + 1
		selectedTags := tags[(index - 1):(index + 1)]

		if err := mgr.CreateServer(addr, selectedTags); err != nil {
			log.Fatal(err)
		}
	}
}
