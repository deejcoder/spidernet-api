package storage

import (
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
