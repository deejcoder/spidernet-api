package server

import (
	sb "github.com/huandu/go-sqlbuilder"
	log "github.com/sirupsen/logrus"
)

/*
	AddServerTags creates many tags and assoicates them to a server,
	this will be transformed into a transaction
	where particular INSERTs will be forced using a fake UPDATE statement
	to return the ID
*/
func (mgr ServerManager) AddServerTags(id int, tags []string) error {
	for _, tag := range tags {
		err := mgr.AddServerTag(id, tag)
		if err != nil {
			log.Warning(err)
		}
	}
	return nil
}

// AddServerTag assoicates a tag with a server
func (mgr ServerManager) AddServerTag(id int, tag string) error {
	err := mgr.CreateTag(tag)
	if err != nil {
		return err
	}

	tag_id, err := mgr.GetTagId(tag)
	if err != nil {
		return err
	}

	iq := sb.Build(`INSERT INTO server_tags (server_id, tag_id) VALUES ($?, $?) ON CONFLICT DO NOTHING`, id, tag_id)
	sql, args := iq.BuildWithFlavor(sb.PostgreSQL)

	_, err = mgr.client.Db.Query(sql, args...)
	return err
}

// CreateTag creates a new tag if it doesn't already exist
func (mgr ServerManager) CreateTag(tag string) error {
	ib := sb.Build(`INSERT INTO tags (tag) VALUES ($?) ON CONFLICT DO NOTHING;`, tag)
	sql, args := ib.BuildWithFlavor(sb.PostgreSQL)
	_, err := mgr.client.Db.Query(sql, args...)
	return err
}

// GetTagId returns the ID of a given tag
func (mgr ServerManager) GetTagId(tag string) (int, error) {
	selectb := sb.PostgreSQL.NewSelectBuilder()
	selectb.Select("id")
	selectb.From("tags")
	selectb.Where(selectb.E("tag", tag))
	selectb.Limit(1)

	sql, args := selectb.Build()
	rows, err := mgr.client.Db.Query(sql, args...)
	if err != nil {
		return 0, err
	}

	var tag_id int
	rows.Next()
	err = rows.Scan(&tag_id)
	return tag_id, err
}
