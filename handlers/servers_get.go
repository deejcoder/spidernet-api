/*
GetServers returns 20 servers starting at START
*/

package handlers

import (
	"net/http"
	"strconv"

	"github.com/deejcoder/spidernet-api/storage"

	"github.com/deejcoder/spidernet-api/helpers"
)

func GetServers(w http.ResponseWriter, r *http.Request) {
	ac := helpers.GetAppContext(r)
	svrmgr := storage.NewServerManager(ac.PostgresInstance.Db)
	response := helpers.NewResponse(w, r)

	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil || start < 0 {
		start = 0
	}

	servers := svrmgr.GetServers(start, 20)
	response.Success("returned some more servers", servers)
}
