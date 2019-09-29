package handlers

import (
	"net/http"
	"strconv"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage/server"
)

func Servers(w http.ResponseWriter, r *http.Request) {
	ac := helpers.GetAppContext(r)
	svrmgr := server.NewServerManager(ac.PostgresInstance)
	response := helpers.NewResponse(w, r)

	// get params & check them
	term := r.URL.Query().Get("term")
	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil {
		response.Error("Invalid Start", helpers.ErrorInvalidParam)
		return
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		response.Error("Invalid Size", helpers.ErrorInvalidParam)
		return
	}

	if size > 40 || size < 1 {
		response.Error("Invalid Start amount", helpers.ErrorInvalidParam)
		return
	}

	servers, err := svrmgr.SearchServers(term, start, size)
	if err != nil {
		response.Error("Search unsuccessful", helpers.ErrorInternalError)
		return
	}

	response.Success("Search successful", servers)
}
