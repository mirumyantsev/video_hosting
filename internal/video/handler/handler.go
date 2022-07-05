package handler

import (
	"github.com/gin-gonic/gin"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/internal/video"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/timedate"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type VideoHandler struct {
	cfg         *config.Config
	useCase     video.VideoUseCase
	logUseCase  logger.LogUseCase
	authUseCase auth.AuthUseCase
	sessUseCase sess.SessUseCase
	userUseCase user.UserUseCase
}

func NewVideoHandler(cfg *config.Config, useCase video.VideoUseCase,
	logUseCase logger.LogUseCase, authUseCase auth.AuthUseCase,
	sessUseCase sess.SessUseCase, userUseCase user.UserUseCase) *VideoHandler {
	return &VideoHandler{
		cfg:         cfg,
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
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputVideo.Url, inputVideo.File) {
		h.logUseCase.Report(ctx, log, msg.ErrorUrlAndFilenameCannotBeEmpty())
		return
	}

	// Assign user ID into info and creation date, create info
	inputVideo.UserId = userId
	inputVideo.CreateDate = log.CreationDate

	if err := h.useCase.CreateVideo(inputVideo); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCreateVideo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoVideoCreated())
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
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsVideoExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckVideoExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorVideoWithRequestedIDIsNotExist())
		return
	}

	gottenVideo, err := h.useCase.GetVideo(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetVideo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotVideo(gottenVideo))
}

func (h *VideoHandler) GetAllVideos(ctx *gin.Context) {
	actPermission := "get_all_videos"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	urlparams := h.userUseCase.ParseURLParams(ctx)

	// Get all infos. If gotten is nothing - send such a message
	gottenVideos, err := h.useCase.GetAllVideos(urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetAllVideos(err))
		return
	}

	if gottenVideos == nil {
		h.logUseCase.Report(ctx, log, msg.InfoNoVideosAvailable())
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotAllVideos(gottenVideos))
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
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsVideoExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckVideoExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorVideoWithRequestedIDIsNotExist())
		return
	}

	// Read input, define ID as requested, partially update info
	inputVideo, err := h.useCase.BindJSONVideo(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	inputVideo.Id = reqId

	if err := h.useCase.PartiallyUpdateVideo(inputVideo); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotPartiallyUpdateVideo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoVideoPartiallyUpdated())
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
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsVideoExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckVideoExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorVideoWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteVideo(reqId); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteVideo(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoVideoDeleted())
}

func (h *VideoHandler) isPermsGranted_getUserId(ctx *gin.Context, log *logger.Log, permission string) (bool, int) {
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
