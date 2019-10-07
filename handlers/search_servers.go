package handlers

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage"
)

func SearchServers(w http.ResponseWriter, r *http.Request) {
	ac := helpers.GetAppContext(r)
	svrmgr := storage.NewServerManager(ac.PostgresInstance.Db)
	response := helpers.NewResponse(w, r)

	term := r.URL.Query().Get("term")
	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil {
		response.Error("start parameter is invalid, this specifies the offset", helpers.ErrorValidationError)
		return
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		response.Error("size parameter is invalid, this specifies the amount of results returned", helpers.ErrorValidationError)
		return
	}

	if size > 40 || size < 1 {
		response.Error("start parameter contains an invalid value", helpers.ErrorValidationError)
		return
	}

	servers, err := svrmgr.SearchServers(term, start, size)
	if err != nil {
		logrus.Error(err)
	}

	response.Success("search was successful", servers)

}
