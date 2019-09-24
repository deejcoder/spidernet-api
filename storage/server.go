package storage

import (
	"strings"
	"time"
)

type Server struct {
	ID           int       `json:"id"`
	Addr         string    `json:"addr"`
	Nick         string    `json:"nick"`
	VotesUp      string    `json:"votes_up"`
	VotesDown    string    `json:"votes_down"`
	Tags         []string  `json:"tags"`
	LastModified time.Time `json:"last_modified"`
	DateAdded    time.Time `json:"date_added"`
}

func (instance PostgresInstance) CreateServer(host string, nick string, tags []string) error {

	tagsString := strings.Join(tags, ",")
	_, err := instance.client.Query(`
		INSERT INTO servers 
		(addr, nick, tags) 
		VALUES ($1, $2, $3)
	`, host, nick, tagsString)

	if err != nil {
		return err
	}
	return nil
}
