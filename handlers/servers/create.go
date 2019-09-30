package servers

import (
	"encoding/json"
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage/server"
)

func Create(w http.ResponseWriter, r *http.Request) {
	ac := helpers.GetAppContext(r)
	smgr := server.NewServerManager(ac.PostgresInstance)
	response := helpers.NewResponse(w, r)

	var server server.ServerWithTags
	err := json.NewDecoder(r.Body).Decode(&server)
	if err != nil {
		response.Error("Invalid server type", helpers.ErrorValidationError)
		return
	}

	err = smgr.UpdateServer(&server.Server)
	if err != nil {
		response.Error("Failed to update server", helpers.ErrorValidationError)
		return
	}

	err = smgr.CreateServerTags(server.Server.ID, server.Tags)
	if err != nil {
		response.Error("Failed to update server tags", helpers.ErrorValidationError)
		return
	}

	response.Success("Update successful", nil)
}
