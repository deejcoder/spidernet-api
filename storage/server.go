package storage

import (
	"context"

	"github.com/deejcoder/spidernet-api/util/config"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Addr      string   `json:"addr"`
	VotesUp   int      `json:"votesUp"`
	VotesDown int      `json:"votesDown"`
	Tags      []string `json:"tags"`
}

// CreateServer inserts a new doc representing a Server
func CreateServer(client *elastic.Client, addr string, tags []string) (bool, error) {
	config := config.GetConfig()
	rootIndex := config.Database.DefaultIndice

	server := Server{
		Addr: addr,
		Tags: tags,
	}

	// OpType=create will avoid overwriting existing docs
	// the id of a record should be the addr
	_, err := client.Index().
		Index(rootIndex).
		Type(TypeServer).
		OpType("create").
		Id(addr).
		BodyJson(server).
		Do(context.Background())

	if err != nil {
		log.Error(err)
		// report back if server exists already
		if elastic.IsConflict(err) {
			return true, nil
		}
		return false, err
	}

	log.WithFields(log.Fields{
		"addr": addr,
	}).Infof("Created a new Server document")
	return false, nil
}

// ServerExists checks if a server with a particular address exists
func ServerExists(client *elastic.Client, addr string) bool {
	config := config.GetConfig()
	rootIndex := config.Database.DefaultIndice

	result, err := client.Get().
		Index(rootIndex).
		Type(TypeServer).
		Id(addr).
		Do(context.Background())

	if err != nil {
		log.Error(err)
		return false
	}

	if result.Found {
		return true
	}
	return false
}
