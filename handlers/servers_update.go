/*
UpdateServer allows some client to update server information,
and update server tags (this includes deleting existing tags, and adding non-existing tags)
*/

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage"
)

func UpdateServer(w http.ResponseWriter, r *http.Request) {
	ac := helpers.GetAppContext(r)
	response := helpers.NewResponse(w, r)
	svrmgr := storage.NewServerManager(ac.PostgresInstance.Db)

	// get the data
	var serverJSON struct {
		Addr string   `json:"addr"`
		Nick string   `json:"nick"`
		Tags []string `json:"tags"`
	}

	err := json.NewDecoder(r.Body).Decode(&serverJSON)
	if err != nil {
		response.Error("server changes rejected!", helpers.ErrorValidationError)
		return
	}

	// update the server
	err = svrmgr.UpdateServer(&storage.Server{
		Addr: serverJSON.Addr,
		Nick: serverJSON.Nick,
	})

	if err != nil {
		response.Error("failed to update the server", helpers.ErrorInternalError)
		return
	}

	response.Success("server updated", nil)

	// update server tags (tags need to be deleted here too!)
	// ...

}
