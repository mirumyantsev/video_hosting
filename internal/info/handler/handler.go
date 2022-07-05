package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/info"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	sconfig "github.com/mirumyantsev/video_hosting/pkg/config_stream"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/timedate"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type InfoHandler struct {
	cfg         *config.Config
	scfg        *sconfig.Config
	useCase     info.InfoUseCase
	logUseCase  logger.LogUseCase
	authUseCase auth.AuthUseCase
	sessUseCase sess.SessUseCase
	userUseCase user.UserUseCase
}

func NewInfoHandler(cfg *config.Config, scfg *sconfig.Config, useCase info.InfoUseCase,
	logUseCase logger.LogUseCase, authUseCase auth.AuthUseCase,
	sessUseCase sess.SessUseCase, userUseCase user.UserUseCase) *InfoHandler {
	return &InfoHandler{
		cfg:         cfg,
		scfg:        scfg,
		useCase:     useCase,
		logUseCase:  logUseCase,
		authUseCase: authUseCase,
		sessUseCase: sessUseCase,
		userUseCase: userUseCase,
	}
}

func (h *InfoHandler) CreateInfo(ctx *gin.Context) {
	actPermission := "post_info"

	log := logger.Init(ctx)

	hasPerms, userId := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read input, check required fields
	inputInfo, err := h.useCase.BindJSONInfo(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputInfo.Stream) {
		h.logUseCase.Report(ctx, log, msg.ErrorStreamCannotBeEmpty())
		return
	}

	// Assign user ID into info and creation date, create info
	inputInfo.UserId = userId
	inputInfo.CreateDate = log.CreationDate

	if err := h.useCase.CreateInfo(inputInfo); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCreateInfo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoInfoCreated())
}

func (h *InfoHandler) GetInfo(ctx *gin.Context) {
	actPermission := "get_info"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check info existence, get info
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsInfoExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckInfoExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorInfoWithRequestedIDIsNotExist())
		return
	}

	gottenInfo, err := h.useCase.GetInfo(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetInfo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotInfo(gottenInfo))
}

func (h *InfoHandler) GetAllInfos(ctx *gin.Context) {
	actPermission := "get_all_infos"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	urlparams := h.userUseCase.ParseURLParams(ctx)

	// Get all infos. If gotten is nothing - send such a message
	gottenInfos, err := h.useCase.GetAllInfos(urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetAllInfos(err))
		return
	}

	if gottenInfos == nil {
		h.logUseCase.Report(ctx, log, msg.InfoNoInfosAvailable())
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotAllInfos(gottenInfos))
}

func (h *InfoHandler) PartiallyUpdateInfo(ctx *gin.Context) {
	actPermission := "patch_info"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check info existence
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsInfoExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckInfoExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorInfoWithRequestedIDIsNotExist())
		return
	}

	// Read input, define ID as requested, partially update info
	inputInfo, err := h.useCase.BindJSONInfo(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	inputInfo.Id = reqId

	if err := h.useCase.PartiallyUpdateInfo(inputInfo); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotPartiallyUpdateInfo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoInfoPartiallyUpdated())
}

func (h *InfoHandler) DeleteInfo(ctx *gin.Context) {
	actPermission := "delete_info"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check info existence, delete info
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsInfoExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckInfoExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorInfoWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteInfo(reqId); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteInfo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoInfoDeleted())
}

func (h *InfoHandler) isPermsGranted_getUserId(ctx *gin.Context, log *logger.Log, permission string) (bool, int) {
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
