package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique; not null"`
	SaltedPassword string `gorm:"not null"`
	IsAdmin        bool
}

type PublicUser struct {
	Username   string `json:"username"`
	IsAdmin    bool   `json:"administrator"`
	Authorized bool   `json:"authorized"`
}

type UserManager struct {
	Db *gorm.DB
}

func NewUserManager(db *gorm.DB) *UserManager {
	return &UserManager{Db: db}
}

func (mgr UserManager) CreateUser(username string, password string, admin bool) error {
	// hash the password
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// create the user
	user := User{
		Username:       username,
		SaltedPassword: string(salt),
		IsAdmin:        admin,
	}

	mgr.Db.FirstOrCreate(&user)
	if err := mgr.Db.Error; err != nil {
		return err
	}
	return nil
}

// AuthorizeUser compares passwords, and returns public user data (PublicUser)
func (mgr UserManager) AuthorizeUser(username string, password string) (PublicUser, error) {
	var puser PublicUser

	// get the user
	var user User
	mgr.Db.First(&user, User{Username: username})

	// compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.SaltedPassword), []byte(password))
	if err != nil {
		return puser, err
	}

	// output public user data
	puser = PublicUser{
		Username:   username,
		IsAdmin:    user.IsAdmin,
		Authorized: true,
	}

	return puser, nil
}

func (mgr UserManager) DeleteUser(username string) error {
	mgr.Db.Delete(&User{}, &User{Username: username})
	if err := mgr.Db.Error; err != nil {
		return err
	}

	return nil
}

func (pu PublicUser) String() string {
	return fmt.Sprintf("User<(%s), admin: %v, authorized: %v>", pu.Username, pu.IsAdmin, pu.Authorized)
}
