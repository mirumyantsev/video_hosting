package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/group"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	perm "github.com/mirumyantsev/video_hosting/internal/permission"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/timedate"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type PermHandler struct {
	cfg          *config.Config
	useCase      perm.PermUseCase
	logUseCase   logger.LogUseCase
	authUseCase  auth.AuthUseCase
	sessUseCase  sess.SessUseCase
	userUseCase  user.UserUseCase
	groupUseCase group.GroupUseCase
}

func NewPermHandler(cfg *config.Config, useCase perm.PermUseCase,
	logUseCase logger.LogUseCase, authUseCase auth.AuthUseCase,
	sessUseCase sess.SessUseCase, userUseCase user.UserUseCase,
	groupUseCase group.GroupUseCase) *PermHandler {
	return &PermHandler{
		cfg:          cfg,
		useCase:      useCase,
		logUseCase:   logUseCase,
		authUseCase:  authUseCase,
		sessUseCase:  sessUseCase,
		userUseCase:  userUseCase,
		groupUseCase: groupUseCase,
	}
}

func (h *PermHandler) GetAllPermissions(ctx *gin.Context) {
	actPermission := "get_all_perms"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	urlparams := h.userUseCase.ParseURLParams(ctx)

	gottenPerms, err := h.useCase.GetAllPermissions(urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetAllPerms(err))
		return
	}

	if gottenPerms == nil {
		h.logUseCase.Report(ctx, log, msg.InfoNoPermsAvailable())
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotAllPerms(gottenPerms))
}

func (h *PermHandler) isPermsGranted_getUserId(ctx *gin.Context, log *logger.Log, permission string) (bool, int) {
	headerToken := h.authUseCase.ReadHeader(ctx)
	if !h.authUseCase.IsTokenExists(headerToken) {
		h.logUseCase.Report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	session, err := h.sessUseCase.GetSessionAndDate(headerToken)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetSessionAndDate(err))
		return false, -1
	}
	if !h.authUseCase.IsSessionExists(session) {
		h.logUseCase.Report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	if timedate.IsDateExpired(session.CreationDate, h.cfg.SessionTTLHours) {
		if err := h.sessUseCase.DeleteSession(headerToken); err != nil {
			h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteSession(err))
			return false, -1
		}
		h.logUseCase.Report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	headerNamepass, err := h.authUseCase.ParseToken(headerToken)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotParseToken(err))
		return false, -1
	}

	gottenUserId, err := h.userUseCase.GetUserId(headerNamepass.Username)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return false, -1
	}
	if gottenUserId < 0 {
		if err := h.sessUseCase.DeleteSession(headerToken); err != nil {
			h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteSession(err))
			return false, -1
		}
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithThisUsernameIsNotExist())
		return false, -1
	}

	log.SessionOwner = headerNamepass.Username

	isSUorStaff := false
	hasPersonalPerm := false
	if isSUorStaff, err = h.userUseCase.IsUserSuperuserOrStaff(headerNamepass.Username); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckSuperuserStaffPermissions(err))
		return false, -1
	}
	if !isSUorStaff {
		if hasPersonalPerm, err = h.userUseCase.IsUserHavePersonalPermission(gottenUserId, permission); err != nil {
			h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckPersonalPermission(err))
			return false, -1
		}
	}

	if !isSUorStaff && !hasPersonalPerm {
		h.logUseCase.Report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	return true, gottenUserId
}
