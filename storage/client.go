package storage

import (
	"context"
	"fmt"

	config "github.com/deejcoder/spidernet-api/util/config"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// Connect connects to a database management system (ElasticSearch) and tests the connection
func Connect() *elastic.Client {

	config := config.GetConfig()

	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	addr := fmt.Sprintf("%s:%d", config.Database.Host, config.Database.Port)
	info, code, err := client.Ping(addr).Do(ctx)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	log.Infof("ElasticSearch returned code %d and version %s\n", code, info.Version.Number)
	return client

}
