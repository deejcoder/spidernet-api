package storage

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/deejcoder/spidernet-api/util/config"
	log "github.com/sirupsen/logrus"
)

func setup() (*config.Config, *PostgresInstance) {
	cfg := config.InitConfig("../")

	// configure & connect to postgres
	psql := NewPostgresInstance()
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

	// search the servers for two tags: community & forum, limit results between 1 & 2
	servers, err := mgr.SearchServers("community forum", 1, 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(servers)

	if len(servers) != 2 {
		log.Panic("Oh dear :( I should get two servers back!")
	}

	// delete server 0
	if err := mgr.DeleteServer(servers[0].ID); err != nil {
		log.Panic(err)
	}

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
		addr := fmt.Sprintf("122.158.3.%d", i)

		// randomly select two tags
		index := rnd.Intn(len(tags)-1) + 1
		selectedTags := tags[(index - 1):(index + 1)]

		if err := mgr.CreateServer(addr, selectedTags); err != nil {
			log.Fatal(err)
		}
	}
}
