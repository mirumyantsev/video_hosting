package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type PermCommon interface {
	GetAllPermissions(urlparams *user.Pagin) (map[int]*Perm, error)
	GetUserPermissions(id int, urlparams *user.Pagin) (*PermIds, error)
	GetGroupPermissions(id int, urlparams *user.Pagin) (*PermIds, error)
}

type PermUseCase interface {
	PermCommon

	SetUserPermissions(id int, permIds *PermIds) error
	DeleteUserPermissions(id int, permIds *PermIds) error

	SetGroupPermissions(id int, permIds *PermIds) error
	DeleteGroupPermissions(id int, permIds *PermIds) error

	BindJSONPermIds(ctx *gin.Context) (*PermIds, error)
	IsRequiredEmpty(permIds *PermIds) bool
	AtoiRequestedId(ctx *gin.Context) (int, error)
}

type PermRepository interface {
	PermCommon

	SetUserPermissions(values string) error
	DeleteUserPermissions(id int, condIds string) error

	SetGroupPermissions(values string) error
	DeleteGroupPermissions(id int, condIds string) error
}
