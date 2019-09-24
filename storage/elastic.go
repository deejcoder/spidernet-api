package storage

import (
	"github.com/olivere/elastic/v7"
)

type ElasticInstance struct{}

// Connect connects to a database management system (ElasticSearch) and tests the connection
func (instance ElasticInstance) Connect() *elastic.Client {
	return nil
}
