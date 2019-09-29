package server

import (
	sql "database/sql"
	"time"

	client "github.com/deejcoder/spidernet-api/storage/client"
	sb "github.com/huandu/go-sqlbuilder"
	log "github.com/sirupsen/logrus"
)

type ServerManager struct {
	client *client.PostgresInstance
}

type Server struct {
	ID           int            `db:"id"`
	Addr         string         `db:"addr"`
	Port         sql.NullInt64  `db:"port"`
	Nick         sql.NullString `db:"nick"`
	VotesUp      string         `db:"votes_up"`
	VotesDown    string         `db:"votes_down"`
	LastModified time.Time      `db:"last_modified"`
	DateAdded    time.Time      `db:"date_added"`
}

var serverStruct = sb.NewStruct(new(Server))

// NewServerManager creates a new manager, to manage the server schema
func NewServerManager(instance *client.PostgresInstance) *ServerManager {
	return &ServerManager{instance}
}

func (mgr ServerManager) CreateServer(host string, tags []string) error {
	iq := sb.Build(`INSERT INTO servers (addr) VALUES ($?) RETURNING id`, host)
	sql, args := iq.BuildWithFlavor(sb.PostgreSQL)

	var id int
	err := mgr.client.Db.QueryRow(sql, args...).Scan(&id)
	if err != nil {
		return err
	}

	// add tags to this server
	err = mgr.CreateServerTags(id, tags)
	return err
}

func (mgr ServerManager) DeleteServer(id int) error {
	err := mgr.client.Delete("servers", "id", id)
	return err
}

func (mgr ServerManager) UpdateServer(server *Server) error {
	err := mgr.client.Update("servers", serverStruct, server)
	return err
}

// SearchServers searches all servers for best matches against tags, and orders by most popular
func (mgr ServerManager) SearchServers(term string, start int, size int) ([]Server, error) {
	var servers []Server
	qb := sb.Build(`
		SELECT 
			servers.*
		FROM servers
		INNER JOIN server_tags
		ON server_tags.server_id = servers.id
		INNER JOIN tags
		ON tags.id = server_tags.tag_id
		ORDER BY 
			similarity(tags.tag, $0) DESC,
			servers.votes_up + servers.votes_down DESC
		OFFSET $1
		LIMIT $2;
	`, term, start, size)

	sql, args := qb.BuildWithFlavor(sb.PostgreSQL)
	hits, err := mgr.client.Db.Query(sql, args...)

	if err != nil {
		return servers, err
	}

	defer hits.Close()
	for hits.Next() {
		server := Server{}
		err := hits.Scan(serverStruct.Addr(&server)...)
		if err != nil {
			log.Error(err)
			continue
		}

		servers = append(servers, server)
	}
	return servers, nil
}
