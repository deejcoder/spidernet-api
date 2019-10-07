/*
DeleteServer will delete some server, given the 'id' param is valid
It will also remove any related tags
*/

package handlers

import (
	"net/http"
	"strconv"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage"
)

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	ac := helpers.GetAppContext(r)
	response := helpers.NewResponse(w, r)
	svrmgr := storage.NewServerManager(ac.PostgresInstance.Db)

	// get provided server id, assure it's unsigned int
	sid, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil {
		response.Error("invalid server id provided", helpers.ErrorInvalidParam)
		return
	}

	err = svrmgr.DeleteServer(uint(sid))
	if err != nil {
		response.Error("failed to delete server with provided id", helpers.ErrorInternalError)
		return
	}

	response.Success("the server has been deleted", nil)
}
