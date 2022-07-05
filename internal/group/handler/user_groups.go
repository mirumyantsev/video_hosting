package handler

import (
	"github.com/gin-gonic/gin"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func (h *GroupHandler) SetUserGroups(ctx *gin.Context) {
	actPermission := "set_user_groups"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, read input, check required fields
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	inputGroupIds, err := h.useCase.BindJSONGroupIds(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsGroupIdsRequiredEmpty(inputGroupIds) {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupIdsCannotBeEmpty())
		return
	}

	// Check user existence, upsert user permissions
	exists, err := h.userUseCase.IsUserExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.SetUserGroups(reqId, inputGroupIds); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotSetUserGroups(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserGroupsSet())
}

func (h *GroupHandler) GetUserGroups(ctx *gin.Context) {
	actPermission := "get_user_groups"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check user existence, get user permissions
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.userUseCase.IsUserExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	urlparams := h.userUseCase.ParseURLParams(ctx)

	gottenGroups, err := h.useCase.GetUserGroups(reqId, urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetUserGroups(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotUserGroups(gottenGroups))
}

func (h *GroupHandler) DeleteUserGroups(ctx *gin.Context) {
	actPermission := "delete_user_groups"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, read input, check required fields
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	inputGroupIds, err := h.useCase.BindJSONGroupIds(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsGroupIdsRequiredEmpty(inputGroupIds) {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupIdsCannotBeEmpty())
		return
	}

	// Check user existence, delete user permissions
	exists, err := h.userUseCase.IsUserExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteUserGroups(reqId, inputGroupIds); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteUserGroups(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserGroupsDeleted())
}
