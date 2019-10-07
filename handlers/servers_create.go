/*
CreateServer creates a new server, and creates any assoicated tags
*/

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
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
		response.Error("server rejected!", helpers.ErrorValidationError)
		return
	}

	// create the server
	server, err := svrmgr.CreateServer(serverJSON.Addr, serverJSON.Nick)
	if err != nil {
		response.Error("we couldn't process the new server", helpers.ErrorInternalError)
		return
	}

	// append any tags to the server
	err = svrmgr.AddServerTags(server, serverJSON.Tags)
	if err != nil {
		response.Success("we created the new server, but we couldn't add the tags", server)
		return
	}

	response.Success("the server was successfully created", server)
}
