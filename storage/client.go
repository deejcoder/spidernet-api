package storage

import (
	"context"
	"fmt"

	config "github.com/deejcoder/spidernet-api/util/config"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

const (
	TypeServer = "server"
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

	addr := fmt.Sprintf("http://%s:%d", config.Database.Host, config.Database.Port)
	info, code, err := client.Ping(addr).Do(ctx)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	log.Infof("ElasticSearch returned code %d and version %s\n", code, info.Version.Number)

	createIndex(client)
	return client

}

// createIndex creates a new main index for us to use for storing server, or user data
func createIndex(client *elastic.Client) {
	config := config.GetConfig()
	index := config.Database.DefaultIndice

	ctx := context.Background()

	// check if index exists
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Panic(err)
	}
	if exists {
		return
	}

	// create if it doesn't exist
	body := `{}`
	result, err := client.CreateIndex(index).BodyString(body).Do(ctx)
	if err != nil {
		log.Panic(err)
	}

	if !result.Acknowledged {
		log.Warning("ElasticSearch has not acknowledged inserting index 'spidernet'")
		return
	}

	log.Infof("Successfully created index=%s", result.Index)
}
