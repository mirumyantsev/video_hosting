package usecase

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/group"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func (u *GroupUseCase) SetUserGroups(id int, permIds *group.GroupIds) error {
	values := ""
	for _, val := range permIds.Ids {
		values += fmt.Sprintf("(%d,%d),", id, val)
	}
	values = values[:len(values)-1]
	return u.groupRepo.SetUserGroups(values)
}

func (u *GroupUseCase) GetUserGroups(id int, urlparams *user.Pagin) (*group.GroupIds, error) {
	return u.groupRepo.GetUserGroups(id, urlparams)
}

func (u *GroupUseCase) DeleteUserGroups(id int, groupIds *group.GroupIds) error {
	condIds := ""
	for _, val := range groupIds.Ids {
		condIds += fmt.Sprintf("%d,", val)
	}
	condIds = condIds[:len(condIds)-1]
	return u.groupRepo.DeleteUserGroups(id, condIds)
}

func (u *GroupUseCase) BindJSONGroupIds(ctx *gin.Context) (*group.GroupIds, error) {
	var groupIds group.GroupIds
	if err := ctx.BindJSON(&groupIds); err != nil {
		return &groupIds, err
	}
	return &groupIds, nil
}

func (u *GroupUseCase) IsGroupIdsRequiredEmpty(groupIds *group.GroupIds) bool {
	if len(groupIds.Ids) == 0 {
		return true
	}
	return false
}
