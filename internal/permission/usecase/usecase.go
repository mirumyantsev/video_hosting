package usecase

import (
	"strconv"

	"github.com/gin-gonic/gin"
	perm "github.com/mirumyantsev/video_hosting/internal/permission"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type PermUseCase struct {
	permRepo perm.PermRepository
}

func NewPermUseCase(permRepo perm.PermRepository) *PermUseCase {
	return &PermUseCase{
		permRepo: permRepo,
	}
}

func (u *PermUseCase) GetAllPermissions(urlparams *user.Pagin) (map[int]*perm.Perm, error) {
	return u.permRepo.GetAllPermissions(urlparams)
}

func (u *PermUseCase) BindJSONPermIds(ctx *gin.Context) (*perm.PermIds, error) {
	var permIds perm.PermIds
	if err := ctx.BindJSON(&permIds); err != nil {
		return &permIds, err
	}
	return &permIds, nil
}

func (u *PermUseCase) IsRequiredEmpty(perms *perm.PermIds) bool {
	if len(perms.Ids) == 0 {
		return true
	}
	return false
}

func (u *PermUseCase) AtoiRequestedId(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return -1, err
	}
	return id, nil
}
