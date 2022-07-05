package handler

import (
	"github.com/gin-gonic/gin"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func (h *PermHandler) SetUserPermissions(ctx *gin.Context) {
	actPermission := "set_user_perms"

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

	inputPermIds, err := h.useCase.BindJSONPermIds(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputPermIds) {
		h.logUseCase.Report(ctx, log, msg.ErrorPermIdsCannotBeEmpty())
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

	if err := h.useCase.SetUserPermissions(reqId, inputPermIds); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotSetUserPerms(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserPermsSet())
}

func (h *PermHandler) GetUserPermissions(ctx *gin.Context) {
	actPermission := "get_user_perms"

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

	gottenPerms, err := h.useCase.GetUserPermissions(reqId, urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetUserPerms(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotUserPerms(gottenPerms))
}

func (h *PermHandler) DeleteUserPermissions(ctx *gin.Context) {
	actPermission := "delete_user_perms"

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

	inputPermIds, err := h.useCase.BindJSONPermIds(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputPermIds) {
		h.logUseCase.Report(ctx, log, msg.ErrorPermIdsCannotBeEmpty())
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

	if err := h.useCase.DeleteUserPermissions(reqId, inputPermIds); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteUserPerms(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserPermsDeleted())
}
