package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deejcoder/spidernet-api/helpers"
	"github.com/deejcoder/spidernet-api/storage"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	ac := helpers.GetAppContext(r)
	usrmgr := storage.NewUserManager(ac.PostgresInstance.Db)
	response := helpers.NewResponse(w, r)

	// get the provided credentials
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		response.Error("the user credentials could not be parsed", helpers.ErrorValidationError)
		return
	}

	// check if the credentials are correct
	usr, err := usrmgr.AuthorizeUser(credentials.Username, credentials.Password)
	if err != nil {
		response.Error("the credentials provided are invalid", helpers.ErrorNotAuthorized)
		return
	}

	// generate a jwt token, which will be sent back as a cookie
	helpers.AuthorizeClient(w)

	// send back the user data
	response.Success("rights granted", usr)
	return
}
