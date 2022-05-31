package handler

import (
	"github.com/gin-gonic/gin"
	msg "github.com/mikerumy/vhosting/internal/messages"
	"github.com/mikerumy/vhosting/pkg/logger"
)

func (h *GroupHandler) SetUserGroups(ctx *gin.Context) {
	log := logger.Init(ctx)

	actPermission := "set_user_groups"

	hasPerms, _ := h.IsPermissionsCheckedGetId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, read input, check required fields
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	inputGroupIds, err := h.useCase.BindJSONGroupIds(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsGroupIdsRequiredEmpty(inputGroupIds) {
		h.report(ctx, log, msg.ErrorGroupIdsCannotBeEmpty())
		return
	}

	// Check user existence, upsert user permissions
	exists, err := h.userUseCase.IsUserExists(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.SetUserGroups(reqId, inputGroupIds); err != nil {
		h.report(ctx, log, msg.ErrorCannotSetUserGroups(err))
		return
	}

	h.report(ctx, log, msg.InfoUserGroupsSet())
}

func (h *GroupHandler) GetUserGroups(ctx *gin.Context) {
	log := logger.Init(ctx)

	actPermission := "get_user_groups"

	hasPerms, _ := h.IsPermissionsCheckedGetId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check user existence, get user permissions
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.userUseCase.IsUserExists(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	gottenGroups, err := h.useCase.GetUserGroups(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotGetUserGroups(err))
		return
	}

	h.report(ctx, log, msg.InfoGotUserGroups(gottenGroups))
}

func (h *GroupHandler) DeleteUserGroups(ctx *gin.Context) {
	log := logger.Init(ctx)

	actPermission := "delete_user_groups"

	hasPerms, _ := h.IsPermissionsCheckedGetId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, read input, check required fields
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	inputGroupIds, err := h.useCase.BindJSONGroupIds(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsGroupIdsRequiredEmpty(inputGroupIds) {
		h.report(ctx, log, msg.ErrorGroupIdsCannotBeEmpty())
		return
	}

	// Check user existence, delete user permissions
	exists, err := h.userUseCase.IsUserExists(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteUserGroups(reqId, inputGroupIds); err != nil {
		h.report(ctx, log, msg.ErrorCannotDeleteUserGroups(err))
		return
	}

	h.report(ctx, log, msg.InfoUserGroupsDeleted())
}
