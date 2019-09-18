package handlers

import (
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
)

// Index index
func Index(w http.ResponseWriter, r *http.Request) {
	response := helpers.NewResponse(w, r)
	response.Success("Authorization validated", nil)
}
