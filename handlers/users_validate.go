package handlers

import (
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
)

func ValidateUser(w http.ResponseWriter, r *http.Request) {
	response := helpers.NewResponse(w, r)
	if validated := helpers.ValidateClient(r); !validated {
		response.Error("invalid token", helpers.ErrorNotAuthorized)
		return
	}
	response.Success("valid token", nil)
}
