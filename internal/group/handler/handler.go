package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/group"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/timedate"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type GroupHandler struct {
	cfg         *config.Config
	useCase     group.GroupUseCase
	logUseCase  logger.LogUseCase
	authUseCase auth.AuthUseCase
	sessUseCase sess.SessUseCase
	userUseCase user.UserUseCase
}

func NewGroupHandler(cfg *config.Config, useCase group.GroupUseCase,
	logUseCase logger.LogUseCase, authUseCase auth.AuthUseCase,
	sessUseCase sess.SessUseCase, userUseCase user.UserUseCase) *GroupHandler {
	return &GroupHandler{
		cfg:         cfg,
		useCase:     useCase,
		logUseCase:  logUseCase,
		authUseCase: authUseCase,
		sessUseCase: sessUseCase,
		userUseCase: userUseCase,
	}
}

func (h *GroupHandler) CreateGroup(ctx *gin.Context) {
	actPermission := "post_group"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read input, check required fields, check if new group already exists
	inputGroup, err := h.useCase.BindJSONGroup(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsGroupRequiredEmpty(inputGroup.Name) {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupNameCannotBeEmpty())
		return
	}

	exists, err := h.useCase.IsGroupExists(inputGroup.Name)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckGroupExistence(err))
		return
	}
	if exists {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupWithEnteredNameIsExist())
		return
	}

	// Create group
	if err := h.useCase.CreateGroup(inputGroup); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCreateGroup(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGroupCreated())
}

func (h *GroupHandler) GetGroup(ctx *gin.Context) {
	actPermission := "get_group"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsGroupExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckGroupExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupWithRequestedIDIsNotExist())
		return
	}

	gottenGroup, err := h.useCase.GetGroup(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetGroup(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotGroup(gottenGroup))
}

func (h *GroupHandler) GetAllGroups(ctx *gin.Context) {
	actPermission := "get_all_groups"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	urlparams := h.userUseCase.ParseURLParams(ctx)

	gottenGroups, err := h.useCase.GetAllGroups(urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetAllGroups(err))
		return
	}

	if gottenGroups == nil {
		h.logUseCase.Report(ctx, log, msg.InfoNoGroupsAvailable())
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotAllGroups(gottenGroups))
}

func (h *GroupHandler) PartiallyUpdateGroup(ctx *gin.Context) {
	actPermission := "patch_group"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check group for existance
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsGroupExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckGroupExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupWithRequestedIDIsNotExist())
		return
	}

	// Read input, define ID as requested, partially update group
	inputGroup, err := h.useCase.BindJSONGroup(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	inputGroup.Id = reqId

	if err := h.useCase.PartiallyUpdateGroup(inputGroup); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotPartiallyUpdateGroup(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGroupPartiallyUpdated())
}

func (h *GroupHandler) DeleteGroup(ctx *gin.Context) {
	actPermission := "delete_group"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsGroupExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckGroupExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorGroupWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteGroup(reqId); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteGroup(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGroupDeleted())
}

func (h *GroupHandler) isPermsGranted_getUserId(ctx *gin.Context, log *logger.Log, permission string) (bool, int) {
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
