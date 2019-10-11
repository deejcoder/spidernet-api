package storage

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestUserSchema(t *testing.T) {
	_, pi := setup()
	mgr := NewUserManager(pi.Db)

	username := "dy1zan"
	password := "thisisatest"
	isadmin := true

	// create a new administrator
	err := mgr.CreateUser(username, password, isadmin)
	if err != nil {
		logrus.Fatal(err)
	}

	// try to login
	_, err = mgr.AuthorizeUser(username, "thisisnotvalid")
	if err == nil {
		logrus.Fatal("There's no error, but there is supposed to be an error!")
	}

	// try to login again
	usr, err := mgr.AuthorizeUser(username, password)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info(usr)
}
