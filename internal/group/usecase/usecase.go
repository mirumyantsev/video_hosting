package usecase

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/group"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type GroupUseCase struct {
	groupRepo group.GroupRepository
}

func NewGroupUseCase(groupRepo group.GroupRepository) *GroupUseCase {
	return &GroupUseCase{
		groupRepo: groupRepo,
	}
}

func (u *GroupUseCase) CreateGroup(grp *group.Group) error {
	return u.groupRepo.CreateGroup(grp)
}

func (u *GroupUseCase) GetGroup(id int) (*group.Group, error) {
	return u.groupRepo.GetGroup(id)
}

func (u *GroupUseCase) GetAllGroups(urlparams *user.Pagin) (map[int]*group.Group, error) {
	return u.groupRepo.GetAllGroups(urlparams)
}

func (u *GroupUseCase) PartiallyUpdateGroup(grp *group.Group) error {
	return u.groupRepo.PartiallyUpdateGroup(grp)
}

func (u *GroupUseCase) DeleteGroup(id int) error {
	return u.groupRepo.DeleteGroup(id)
}

func (u *GroupUseCase) IsGroupExists(idOrName interface{}) (bool, error) {
	exists, err := u.groupRepo.IsGroupExists(idOrName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *GroupUseCase) BindJSONGroup(ctx *gin.Context) (*group.Group, error) {
	var grp group.Group
	if err := ctx.BindJSON(&grp); err != nil {
		return &grp, err
	}
	return &grp, nil
}

func (u *GroupUseCase) IsGroupRequiredEmpty(name string) bool {
	if name == "" {
		return true
	}
	return false
}

func (u *GroupUseCase) AtoiRequestedId(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return -1, err
	}
	return id, nil
}
