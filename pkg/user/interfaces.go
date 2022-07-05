package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
)

type UserCommon interface {
	CreateUser(usr *User) error
	GetUser(id int) (*User, error)
	GetAllUsers(urlparams *Pagin) (map[int]*User, error)
	UpdateUserPassword(namepass *auth.Namepass) error
	PartiallyUpdateUser(usr *User) error
	DeleteUser(id int) error

	IsUserSuperuserOrStaff(username string) (bool, error)
	IsUserHavePersonalPermission(userId int, userPerm string) (bool, error)
	IsUserExists(idOrUsername interface{}) (bool, error)
	GetUserId(username string) (int, error)
}

type UserUseCase interface {
	UserCommon

	BindJSONUser(ctx *gin.Context) (*User, error)
	IsRequiredEmpty(username, password string) bool
	AtoiRequestedId(ctx *gin.Context) (int, error)
	ParseURLParams(ctx *gin.Context) *Pagin
}

type UserRepository interface {
	UserCommon
}
