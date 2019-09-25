package storage

import (
	sql "database/sql"
	"time"

	sb "github.com/huandu/go-sqlbuilder"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type ServerManager struct {
	client *PostgresInstance
}

type Server struct {
	ID           int            `db:"id"`
	Addr         string         `db:"addr"`
	Port         sql.NullInt64  `db:"port"`
	Nick         sql.NullString `db:"nick"`
	VotesUp      string         `db:"votes_up"`
	VotesDown    string         `db:"votes_down"`
	Tags         pq.StringArray `db:"tags"`
	LastModified time.Time      `db:"last_modified"`
	DateAdded    time.Time      `db:"date_added"`
}

var serverStruct = sb.NewStruct(new(Server))

// NewServerManager creates a new manager, to manage the server schema
func NewServerManager(instance *PostgresInstance) *ServerManager {
	return &ServerManager{instance}
}

func (mgr ServerManager) CreateServer(host string, tags []string) error {
	ib := sb.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("servers")
	ib.Cols("addr", "tags")
	ib.Values(host, pq.Array(tags))

	// execute
	sql, args := ib.Build()
	_, err := mgr.client.db.Query(sql, args...)
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

// SearchServers matches the {term} with server tags, and returns back {size} tuples starting at {start}
func (mgr ServerManager) SearchServers(term string, start int, size int) ([]Server, error) {
	// build the query
	searchb := sb.Build(
		`
			SELECT *
			FROM servers s
			WHERE s.tags && tsvector_to_array(to_tsvector($?))
			ORDER BY s.votes_up + s.votes_down
			LIMIT $?
			OFFSET $?
		`,
		term, size, start,
	)

	// execute the query
	sql, args := searchb.BuildWithFlavor(sb.PostgreSQL)
	rows, err := mgr.client.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	// build array containing results
	var servers []Server
	for rows.Next() {

		s := Server{}
		err := rows.Scan(serverStruct.Addr(&s)...)
		if err != nil {
			log.Errorf("Excluded result from server search results, error=%s", err)
			continue
		}

		servers = append(servers, s)
	}
	return servers, nil
}
