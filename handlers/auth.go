package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	helpers.AuthorizeClient(w)
	conf := helpers.GetAppContext(r).Config
	response := helpers.NewResponse(w, r)

	// this is here temp during development until users are developed
	var body struct {
		Passphase string `json:"secret_key"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		response.Error("internal error", helpers.ErrorInternalError)
		return
	}

	if body.Passphase == conf.Keys.ApiLogin {
		helpers.AuthorizeClient(w)
		response.Success("rights granted", nil)
		return
	}

	response.Error("invalid passphase", helpers.ErrorNotAuthorized)
}

func Validate(w http.ResponseWriter, r *http.Request) {
	response := helpers.NewResponse(w, r)
	if validated := helpers.ValidateClient(r); !validated {
		response.Error("invalid token", helpers.ErrorNotAuthorized)
		return
	}
	response.Success("valid token", nil)
}
