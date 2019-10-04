package storage

import (
	"testing"

	"github.com/deejcoder/spidernet-api/util/config"
	"github.com/sirupsen/logrus"
)

func setup() (*config.Config, *PostgresInstance) {
	conf := config.InitConfig("../")
	pi := NewPostgresInstance()
	err := pi.Connect()
	if err != nil {
		logrus.Fatal(err)
	}

	return conf, pi
}

func TestServerSchema(t *testing.T) {
	_, pi := setup()
	mgr := NewServerManager(pi.Db)

	// create some server, should return an error if one already exists
	server, err := mgr.CreateServer("192.168.1.4", "Heroes & Generals")
	if err != nil {
		logrus.Trace(err)
	}

	// create some tags for the server
	tags := []string{
		"Debian",
		"TCP",
		"Game",
	}
	err = mgr.AddServerTags(server, tags)
	if err != nil {
		logrus.Trace(err)
	}

	// finally search the servers for the new server
	servers := mgr.SearchServers("deb", 0, 10)
	for _, server := range servers {
		logrus.Info(server)
		if server.Addr == mgr.GetServerByAddr("192.168.1.4").Addr {

			// finally delete the server, and we're done
			mgr.DeleteServer(server.ID)
			return
		}
	}
	logrus.Panic("test failed!! D:")
}