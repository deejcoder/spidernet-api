package server

import (
	sql "database/sql"
	"time"

	client "github.com/deejcoder/spidernet-api/storage/client"
	sb "github.com/huandu/go-sqlbuilder"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type ServerManager struct {
	client *client.PostgresInstance
}

type Server struct {
	ID           int            `db:"id" json:"id"`
	Addr         string         `db:"addr" json:"addr"`
	Port         sql.NullInt64  `db:"port" json:"port"`
	Nick         sql.NullString `db:"nick" json:"nick"`
	VotesUp      string         `db:"votes_up" json:"votes_up"`
	VotesDown    string         `db:"votes_down" json:"votes_down"`
	LastModified time.Time      `db:"last_modified" json:"last_modified"`
	DateAdded    time.Time      `db:"date_added" json:"date_added"`
}

type ServerWithTags struct {
	Server Server   `json:"server"`
	Tags   []string `db:"tags" json:"tags"`
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
func (mgr ServerManager) SearchServers(term string, start int, size int) ([]ServerWithTags, error) {
	qb := sb.Build(`
		WITH 
			get_server_tags AS (
				SELECT server_id, tag
				FROM tags
				INNER JOIN server_tags
				ON tags.id = server_tags.tag_id
			),
			tags_to_array AS (
				SELECT server_id, ARRAY_AGG(tag) AS tags
				FROM get_server_tags
				GROUP BY server_id
			)
		
		SELECT servers.*, tags_to_array.tags
		FROM servers
		INNER JOIN server_tags
		ON server_tags.server_id = servers.id
		INNER JOIN tags
		ON tags.id = server_tags.tag_id
		INNER JOIN tags_to_array
		ON tags_to_array.server_id = servers.id
		ORDER BY 
			similarity(tags.tag, $0) DESC,
			servers.votes_up + servers.votes_down DESC
		OFFSET $1
		LIMIT $2
	`, term, start, size)

	// build the query
	sql, args := qb.BuildWithFlavor(sb.PostgreSQL)
	hits, err := mgr.client.Db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer hits.Close()
	servers := mgr.DecodeServers(hits)
	return servers, nil
}

func (mgr ServerManager) DecodeServers(hits *sql.Rows) []ServerWithTags {
	var servers []ServerWithTags
	for hits.Next() {
		server := ServerWithTags{}

		stct := append(
			serverStruct.Addr(&server.Server),
			(*pq.StringArray)(&server.Tags),
		)
		err := hits.Scan(stct...)
		if err != nil {
			log.Error(err)
			continue
		}

		servers = append(servers, server)
	}
	return servers
}
