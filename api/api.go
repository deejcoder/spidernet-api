package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage/client"
	"github.com/deejcoder/spidernet-api/util/config"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

func configure(ac *helpers.AppContext) *http.Server {

	config := config.GetConfig()

	log.Info(config.API.AllowedOrigins)
	cors := handlers.CORS(
		handlers.AllowedOrigins(config.API.AllowedOrigins),
		handlers.AllowedMethods(config.API.AllowedMethods),
		handlers.AllowedHeaders(config.API.AllowedHeaders),
	)

	router := BuildRouter()

	// enable csrf tokens
	csrfMiddleware := csrf.Protect(
		[]byte(config.Keys.CSRFKey),
		csrf.Secure(config.API.UsingHttps),
	)

	handler := cors(router)
	handler = helpers.AttachCSRFToken(handler)
	handler = csrfMiddleware(handler)
	handler = helpers.AttachAppContext(ac, handler)
	handler = helpers.AttachLogging(handler)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", ac.Config.API.Port),
		Handler: handler,
	}

	return s
}

// Start starts the webserver, terminates on request
func Start(ctx context.Context) {
	conf := config.GetConfig()

	instance := client.NewPostgresInstance()
	if err := instance.Connect(); err != nil {
		log.Fatal(err)
	}

	appContext := helpers.AppContext{
		PostgresInstance: instance,
		Config:           conf,
	}

	server := configure(&appContext)

	// listen for interupt signal to close server
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error(err)
		}
		close(done)
	}()

	log.Infof("Starting REST api on http://localhost:%d", appContext.Config.API.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}

	<-done
}
