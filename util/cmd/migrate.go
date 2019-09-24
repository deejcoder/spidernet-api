package cmd

import (
	"github.com/deejcoder/spidernet-api/storage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "migrates the PostgreSQL database with changes in the migrations folder",
	RunE: func(cmd *cobra.Command, args []string) error {

		instance := storage.NewPostgresInstance()

		db, err := instance.Connect()
		if err != nil {
			log.Fatal(err)
		}

		if err := instance.Migrate(db); err != nil {
			log.Fatal(err)
		}
		return nil
	},
}