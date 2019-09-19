package storage

import (
	"context"

	"github.com/deejcoder/spidernet-api/util/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Addr      string   `json:"addr"`
	VotesUp   int      `json:"votesUp"`
	VotesDown int      `json:"votesDown"`
	Tags      []string `json:"tags"`
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

// GetServer returns a server which exists
func GetServer(client *elastic.Client, addr string) *Server {
	config := config.GetConfig()
	rootIndex := config.Database.DefaultIndice

	result, err := client.Get().
		Index(rootIndex).
		Type(TypeServer).
		Id(addr).
		Do(context.Background())

	if err != nil {
		log.Error(err)
		return nil
	}

	if result.Found {

		// deserialize the json object
		var server *Server
		if err := json.Unmarshal(result.Source, &server); err != nil {
			log.Error(err)
			return nil
		}
		return server
	}
	return nil
}

func GetServers(client *elastic.Client) []*Server {
	config := config.GetConfig()
	rootIndex := config.Database.DefaultIndice

	result, err := client.Search().
		Index(rootIndex).
		Type(TypeServer).
		Do(context.Background())

	if err != nil {
		log.Error(err)
		return nil
	}

	return decodeServers(result)
}

// decodeServers accepts a Search Result and decodes it into []*Server
func decodeServers(result *elastic.SearchResult) []*Server {
	if result == nil || result.TotalHits() == 0 {
		return nil
	}

	var servers []*Server
	for _, hit := range result.Hits.Hits {
		server := new(Server)

		// deserialize each server and append it to servers
		if err := json.Unmarshal(hit.Source, server); err != nil {
			log.Error(err)
			return nil
		}
		servers = append(servers, server)
	}
	return servers
}
