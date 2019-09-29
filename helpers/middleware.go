/*
middleware provides additional functionality to requests
AttachAppContext provides access to the database client and configuration, from requests
*/

package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/deejcoder/spidernet-api/storage/client"
	"github.com/deejcoder/spidernet-api/util/config"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
)

// AppContext allows access to shared app data within handlers
type AppContext struct {
	PostgresInstance *client.PostgresInstance
	Config           *config.Config
}

type key struct{}

// GetAppContext returns the AppContext from a given request for access within Handlers
func GetAppContext(r *http.Request) *AppContext {
	ac, _ := r.Context().Value(key{}).(*AppContext)
	return ac
}

// AttachAppContext attaches an AppContext to a request
func AttachAppContext(appContext *AppContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), key{}, appContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AttachCSRFToken attaches a CSRF token to the header (X-CSRF-Token)
func AttachCSRFToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		next.ServeHTTP(w, r)
	})
}

func AttachLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// logs will be printed at the end of the request
		defer func() {
			log.WithFields(log.Fields{
				"remote":   r.RemoteAddr,
				"duration": time.Since(startTime),
			}).Infof("%s %s", r.Method, r.URL.RequestURI())
		}()

		next.ServeHTTP(w, r)
	})
}
