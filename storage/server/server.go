package server

import (
	sql "database/sql"
	"time"

	client "github.com/deejcoder/spidernet-api/storage/client"
	sb "github.com/huandu/go-sqlbuilder"
	"github.com/lib/pq"
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
	Tags         pq.StringArray `db:"tags"`
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
