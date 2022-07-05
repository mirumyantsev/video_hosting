package handler

import (
	"github.com/gin-gonic/gin"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/timedate"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type UserHandler struct {
	cfg         *config.Config
	useCase     user.UserUseCase
	logUseCase  logger.LogUseCase
	authUseCase auth.AuthUseCase
	sessUseCase sess.SessUseCase
}

func NewUserHandler(cfg *config.Config, useCase user.UserUseCase, logUseCase logger.LogUseCase,
	authUseCase auth.AuthUseCase, sessUseCase sess.SessUseCase) *UserHandler {
	return &UserHandler{
		cfg:         cfg,
		useCase:     useCase,
		logUseCase:  logUseCase,
		authUseCase: authUseCase,
		sessUseCase: sessUseCase,
	}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	actPermission := "post_user"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read input, check required fields, check user existence
	inputUser, err := h.useCase.BindJSONUser(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputUser.Username, inputUser.PasswordHash) {
		h.logUseCase.Report(ctx, log, msg.ErrorUsernameAndPasswordCannotBeEmpty())
		return
	}

	exists, err := h.useCase.IsUserExists(inputUser.Username)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithEnteredUsernameIsExist())
		return
	}

	// Assign user creation time, create user
	inputUser.JoiningDate = log.CreationDate

	if err := h.useCase.CreateUser(inputUser); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCreateUser(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserCreated())
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	actPermission := "get_user"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check user existence, get user
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsUserExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	gottenUser, err := h.useCase.GetUser(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetUser(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotUser(gottenUser))
}

func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	actPermission := "get_all_users"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	urlparams := h.useCase.ParseURLParams(ctx)

	// Get all users. If gotten is nothing - send such a message
	gottenUsers, err := h.useCase.GetAllUsers(urlparams)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetAllUsers(err))
		return
	}

	if gottenUsers == nil {
		h.logUseCase.Report(ctx, log, msg.InfoNoUsersAvailable())
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoGotAllUsers(gottenUsers))
}

func (h *UserHandler) UpdateUserPassword(ctx *gin.Context) {
	actPermission := "post_user_pass"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	inputNamepass, err := h.authUseCase.BindJSONNamepass(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputNamepass.Username, inputNamepass.PasswordHash) {
		h.logUseCase.Report(ctx, log, msg.ErrorUsernameAndPasswordCannotBeEmpty())
		return
	}

	exists, err := h.useCase.IsUserExists(inputNamepass.Username)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithEnteredUsernameIsNotExist())
		return
	}

	if err := h.useCase.UpdateUserPassword(inputNamepass); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotPartiallyUpdateUser(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserPasswordChanged())
}

func (h *UserHandler) PartiallyUpdateUser(ctx *gin.Context) {
	actPermission := "patch_user"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check user existence
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsUserExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	// Read input, define ID as requested, partially update user
	inputUser, err := h.useCase.BindJSONUser(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	inputUser.Id = reqId

	if err := h.useCase.PartiallyUpdateUser(inputUser); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotPartiallyUpdateUser(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserPartiallyUpdated())
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	actPermission := "delete_user"

	log := logger.Init(ctx)

	hasPerms, _ := h.isPermsGranted_getUserId(ctx, log, actPermission)
	if !hasPerms {
		return
	}

	// Read requested ID, check user existence, delete user
	reqId, err := h.useCase.AtoiRequestedId(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotConvertRequestedIDToTypeInt(err))
		return
	}

	exists, err := h.useCase.IsUserExists(reqId)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithRequestedIDIsNotExist())
		return
	}

	if err := h.useCase.DeleteUser(reqId); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteUser(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoUserDeleted())
}

func (h *UserHandler) isPermsGranted_getUserId(ctx *gin.Context, log *logger.Log, permission string) (bool, int) {
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

	gottenUserId, err := h.useCase.GetUserId(headerNamepass.Username)
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
	if isSUorStaff, err = h.useCase.IsUserSuperuserOrStaff(headerNamepass.Username); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckSuperuserStaffPermissions(err))
		return false, -1
	}
	if !isSUorStaff {
		if hasPersonalPerm, err = h.useCase.IsUserHavePersonalPermission(gottenUserId, permission); err != nil {
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
