package handler

import (
	"github.com/gin-gonic/gin"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func (h *PermHandler) SetGroupPermissions(ctx *gin.Context) {
	actPermission := "set_group_perms"

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

	// Check group existence, upsert group permissions
	exists, err := h.groupUseCase.IsGroupExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckGroupExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.SetGroupPermissions(reqId, inputPermIds); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotSetGroupPerms(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGroupPermsSet())
}

func (h *PermHandler) GetGroupPermissions(ctx *gin.Context) {
	actPermission := "get_group_perms"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check group existence, get group permissions
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.groupUseCase.IsGroupExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckGroupExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupWithRequestedIDIsNotExist())
		return
	}

	urlparams := h.userUseCase.ParseURLParams(ctx)

	gottenPerms, err := h.useCase.GetGroupPermissions(reqId, urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetGroupPerms(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotGroupPerms(gottenPerms))
}

func (h *PermHandler) DeleteGroupPermissions(ctx *gin.Context) {
	actPermission := "delete_group_perms"

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

	// Check group existence, delete group permissions
	exists, err := h.groupUseCase.IsGroupExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckGroupExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteGroupPermissions(reqId, inputPermIds); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteGroupPerms(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGroupPermsDeleted())
}
