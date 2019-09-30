/* server_tags
Tags are assoicated with servers.
Contract: A tag will exist if a server exists
If a server no longer exists, a tag must be removed
*/

package server

import (
	sql "database/sql"

	sb "github.com/huandu/go-sqlbuilder"
	log "github.com/sirupsen/logrus"
)

func (mgr ServerManager) CreateServerTags(sid int, tags []string) error {
	tx, err := mgr.client.Db.Begin()
	if err != nil {
		return err
	}

	// add a tag to the transaction
	for _, tag := range tags {
		if err := mgr.CreateServerTag(tx, sid, tag); err != nil {
			// log the error & cancel transaction
			log.Fatal(err)
			err := tx.Rollback()
			return err
		}
	}

	// commit the changes
	err = tx.Commit()
	return err
}

func (mgr ServerManager) CreateServerTag(tx *sql.Tx, sid int, tag string) error {
	_, err := tx.Exec(`
		INSERT INTO
		tags (tag)
		VALUES ($1)
		ON CONFLICT (tag) DO NOTHING;
	`, tag)

	if err != nil {
		log.Errorf("Error inserting new tag: %s", err)
		return err
	}

	_, err = tx.Exec(`
		WITH 
		relations AS (
			SELECT *, $1::INT AS add_to_server_id FROM tags
			LEFT JOIN server_tags
			ON $1::INT = server_tags.server_id
			WHERE server_tags.server_id IS NULL AND tags.tag = $2
		)

		INSERT INTO server_tags
		(server_id, tag_id)
		(
			SELECT add_to_server_id, id
			FROM relations
		)
	`, sid, tag)
	return err
}

func (mgr ServerManager) DeleteServerTags(sid int) error {
	delb := sb.PostgreSQL.NewDeleteBuilder()
	delb.DeleteFrom("server_tags")
	delb.Where(delb.E("server_id", sid))
	sql, args := delb.Build()

	if _, err := mgr.client.Db.Query(sql, args...); err != nil {
		return err
	}
	return nil
}
