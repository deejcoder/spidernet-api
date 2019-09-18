package handlers

import (
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	helpers.AuthorizeClient(w)

	response := helpers.NewResponse(w, r)
	response.Success("Rights authorized", nil)
}
