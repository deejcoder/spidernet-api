package cmd

import (
	"os"

	"github.com/deejcoder/spidernet-api/util/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spidernet",
	Short: "spidernet provides a RESTful service to manage and view a collection of servers.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute executes
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	config.InitConfig(".")
}
