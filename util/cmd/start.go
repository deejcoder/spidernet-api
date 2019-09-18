package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/deejcoder/spidernet-api/api"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCommand)
}

var serveCommand = &cobra.Command{
	Use:   "start",
	Short: "starts the service",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch

			log.Info("signal caught. shutting down...")
			cancel()
		}()

		// cancel when api exits e.g unexpectedly
		defer cancel()
		api.Start(ctx)
		return nil
	},
}
