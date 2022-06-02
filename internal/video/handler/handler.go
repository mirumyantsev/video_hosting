package handler

import (
	"github.com/gin-gonic/gin"
	lg "github.com/mikerumy/vhosting/internal/logging"
	msg "github.com/mikerumy/vhosting/internal/messages"
	sess "github.com/mikerumy/vhosting/internal/session"
	"github.com/mikerumy/vhosting/internal/video"
	"github.com/mikerumy/vhosting/pkg/auth"
	"github.com/mikerumy/vhosting/pkg/logger"
	"github.com/mikerumy/vhosting/pkg/responder"
	"github.com/mikerumy/vhosting/pkg/timedate"
	"github.com/mikerumy/vhosting/pkg/user"
)

type VideoHandler struct {
	useCase     video.VideoUseCase
	logUseCase  lg.LogUseCase
	authUseCase auth.AuthUseCase
	sessUseCase sess.SessUseCase
	userUseCase user.UserUseCase
}

func NewVideoHandler(useCase video.VideoUseCase, logUseCase lg.LogUseCase, authUseCase auth.AuthUseCase,
	sessUseCase sess.SessUseCase, userUseCase user.UserUseCase) *VideoHandler {
	return &VideoHandler{
		useCase:     useCase,
		logUseCase:  logUseCase,
		authUseCase: authUseCase,
		sessUseCase: sessUseCase,
		userUseCase: userUseCase,
	}
}

func (h *VideoHandler) CreateVideo(ctx *gin.Context) {
	actPermission := "post_video"

	log := logger.Init(ctx)

	hasPerms, userId := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read input, check required fields
	inputVideo, err := h.useCase.BindJSONVideo(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputVideo.Url, inputVideo.Filename) {
		h.report(ctx, log, msg.ErrorUrlAndFilenameCannotBeEmpty())
		return
	}

	// Assign user ID into info and creation date, create info
	inputVideo.UserId = userId
	inputVideo.CreationDate = log.CreationDate

	if err := h.useCase.CreateVideo(inputVideo); err != nil {
		h.report(ctx, log, msg.ErrorCannotCreateVideo(err))
		return
	}

	h.report(ctx, log, msg.InfoVideoCreated())
}

func (h *VideoHandler) GetVideo(ctx *gin.Context) {
	actPermission := "get_video"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check info existence, get info
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsVideoExists(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckVideoExistence(err))
		return
	}
	if !exists {
		h.report(ctx, log, msg.ErrorVideoWithRequestedIDIsNotExist())
		return
	}

	gottenVideo, err := h.useCase.GetVideo(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotGetVideo(err))
		return
	}

	h.report(ctx, log, msg.InfoGotVideo(gottenVideo))
}

func (h *VideoHandler) GetAllVideos(ctx *gin.Context) {
	actPermission := "get_all_videos"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Get all infos. If gotten is nothing - send such a message
	gottenVideos, err := h.useCase.GetAllVideos()
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotGetAllVideos(err))
		return
	}

	if gottenVideos == nil {
		h.report(ctx, log, msg.InfoNoVideosAvailable())
		return
	}

	h.report(ctx, log, msg.InfoGotAllVideos(gottenVideos))
}

func (h *VideoHandler) PartiallyUpdateVideo(ctx *gin.Context) {
	actPermission := "patch_video"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check info existence
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsVideoExists(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckVideoExistence(err))
		return
	}
	if !exists {
		h.report(ctx, log, msg.ErrorVideoWithRequestedIDIsNotExist())
		return
	}

	// Read input, define ID as requested, partially update info
	inputVideo, err := h.useCase.BindJSONVideo(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	inputVideo.Id = reqId

	if err := h.useCase.PartiallyUpdateVideo(inputVideo); err != nil {
		h.report(ctx, log, msg.ErrorCannotPartiallyUpdateVideo(err))
		return
	}

	h.report(ctx, log, msg.InfoVideoPartiallyUpdated())
}

func (h *VideoHandler) DeleteVideo(ctx *gin.Context) {
	actPermission := "delete_video"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check info existence, delete info
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsVideoExists(reqId)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckVideoExistence(err))
		return
	}
	if !exists {
		h.report(ctx, log, msg.ErrorVideoWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteVideo(reqId); err != nil {
		h.report(ctx, log, msg.ErrorCannotDeleteVideo(err))
		return
	}

	h.report(ctx, log, msg.InfoVideoDeleted())
}

func (h *VideoHandler) report(ctx *gin.Context, log *lg.Log, messageLog *lg.Log) {
	logger.Complete(log, messageLog)
	responder.Response(ctx, log)
	if err := h.logUseCase.CreateLogRecord(log); err != nil {
		logger.Complete(log, msg.ErrorCannotDoLogging(err))
		responder.Response(ctx, log)
	}
	logger.Print(log)
}

func (h *VideoHandler) isPermsGranted_getUserId(ctx *gin.Context, log *lg.Log, permission string) (bool, int) {
	headerToken := h.authUseCase.ReadHeader(ctx)
	if !h.authUseCase.IsTokenExists(headerToken) {
		h.report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	session, err := h.sessUseCase.GetSessionAndDate(headerToken)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotGetSessionAndDate(err))
		return false, -1
	}
	if !h.authUseCase.IsSessionExists(session) {
		h.report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	if timedate.IsDateExpired(session.CreationDate) {
		if err := h.sessUseCase.DeleteSession(headerToken); err != nil {
			h.report(ctx, log, msg.ErrorCannotDeleteSession(err))
			return false, -1
		}
		h.report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	headerNamepass, err := h.authUseCase.ParseToken(headerToken)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotParseToken(err))
		return false, -1
	}

	gottenUserId, err := h.userUseCase.GetUserId(headerNamepass.Username)
	if err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return false, -1
	}
	if gottenUserId < 0 {
		if err := h.sessUseCase.DeleteSession(headerToken); err != nil {
			h.report(ctx, log, msg.ErrorCannotDeleteSession(err))
			return false, -1
		}
		h.report(ctx, log, msg.ErrorUserWithThisUsernameIsNotExist())
		return false, -1
	}

	log.SessionOwner = headerNamepass.Username

	isSUorStaff := false
	hasPersonalPerm := false
	if isSUorStaff, err = h.userUseCase.IsUserSuperuserOrStaff(headerNamepass.Username); err != nil {
		h.report(ctx, log, msg.ErrorCannotCheckSuperuserStaffPermissions(err))
		return false, -1
	}
	if !isSUorStaff {
		if hasPersonalPerm, err = h.userUseCase.IsUserHavePersonalPermission(gottenUserId, permission); err != nil {
			h.report(ctx, log, msg.ErrorCannotCheckPersonalPermission(err))
			return false, -1
		}
	}

	if !isSUorStaff && !hasPersonalPerm {
		h.report(ctx, log, msg.ErrorYouHaveNotEnoughPermissions())
		return false, -1
	}

	return true, gottenUserId
}