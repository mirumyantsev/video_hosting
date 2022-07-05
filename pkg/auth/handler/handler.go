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

type AuthHandler struct {
	cfg         *config.Config
	useCase     auth.AuthUseCase
	userUseCase user.UserUseCase
	sessUseCase sess.SessUseCase
	logUseCase  logger.LogUseCase
}

func NewAuthHandler(cfg *config.Config, useCase auth.AuthUseCase,
	userUseCase user.UserUseCase, sessUseCase sess.SessUseCase,
	logUseCase logger.LogUseCase) *AuthHandler {
	return &AuthHandler{
		cfg:         cfg,
		useCase:     useCase,
		userUseCase: userUseCase,
		sessUseCase: sessUseCase,
		logUseCase:  logUseCase,
	}
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	log := logger.Init(ctx)

	headerToken := h.useCase.ReadHeader(ctx)
	if h.useCase.IsTokenExists(headerToken) {
		if err := h.sessUseCase.DeleteSession(headerToken); err != nil {
			h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteSession(err))
			return
		}
	}

	inputNamepass, err := h.useCase.BindJSONNamepass(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.userUseCase.IsRequiredEmpty(inputNamepass.Username, inputNamepass.PasswordHash) {
		h.logUseCase.Report(ctx, log, msg.ErrorUsernameAndPasswordCannotBeEmpty())
		return
	}

	exists, err := h.useCase.IsUsernameAndPasswordExists(inputNamepass.Username, inputNamepass.PasswordHash)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithEnteredUsernameOrPasswordIsNotExist())
		return
	}

	log.SessionOwner = inputNamepass.Username

	newToken, err := h.useCase.GenerateToken(inputNamepass)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGenerateToken(err))
		return
	}

	if err := h.sessUseCase.CreateSession(ctx, inputNamepass.Username, newToken, log.CreationDate); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCreateSession(err))
		return
	}

	h.logUseCase.ReportWithToken(ctx, log, msg.InfoYouHaveSuccessfullySignedIn(), newToken)
}

func (h *AuthHandler) ChangePassword(ctx *gin.Context) {
	log := logger.Init(ctx)

	session, err := h.getValidSessionAndDeleteSession(ctx, log)
	if err != nil {
		return
	}
	if session == nil {
		h.logUseCase.Report(ctx, log, msg.ErrorYouMustBeSignedInForChangingPassword())
		return
	}

	sessionNamepass, err := h.useCase.ParseToken(session.Content)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotParseToken(err))
		return
	}

	inputNamepass, err := h.useCase.BindJSONNamepass(ctx)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotBindInputData(err))
		return
	}

	if h.useCase.IsRequiredEmpty(inputNamepass) {
		h.logUseCase.Report(ctx, log, msg.ErrorPasswordCannotBeEmpty())
		return
	}

	inputNamepass.Username = sessionNamepass.Username

	exists, err := h.userUseCase.IsUserExists(inputNamepass.Username)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithSuchUsernameOrPasswordIsNotExist())
		return
	}

	log.SessionOwner = sessionNamepass.Username

	if err := h.useCase.UpdateUserPassword(inputNamepass); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotUpdateUserPassword(err))
		return
	}

	h.logUseCase.Report(ctx, log, msg.InfoYouHaveSuccessfullyChangedPassword())
}

func (h *AuthHandler) SignOut(ctx *gin.Context) {
	log := logger.Init(ctx)

	session, err := h.getValidSessionAndDeleteSession(ctx, log)
	if err != nil {
		return
	}
	if session == nil {
		h.logUseCase.Report(ctx, log, msg.ErrorYouMustBeSignedInForSigningOut())
		return
	}

	sessionNamepass, err := h.useCase.ParseToken(session.Content)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotParseToken(err))
		return
	}

	exists, err := h.userUseCase.IsUserExists(sessionNamepass.Username)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotCheckUserExistence(err))
		return
	}
	if !exists {
		h.logUseCase.Report(ctx, log, msg.ErrorUserWithThisUsernameIsNotExist())
		return
	}

	log.SessionOwner = sessionNamepass.Username

	h.logUseCase.Report(ctx, log, msg.InfoYouHaveSuccessfullySignedOut())
}

func (h *AuthHandler) getValidSessionAndDeleteSession(ctx *gin.Context, log *logger.Log) (*sess.Session, error) {
	headerToken := h.useCase.ReadHeader(ctx)
	if !h.useCase.IsTokenExists(headerToken) {
		return nil, nil
	}

	session, err := h.sessUseCase.GetSessionAndDate(headerToken)
	if err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotGetSessionAndDate(err))
		return nil, err
	}
	if !h.useCase.IsSessionExists(session) {
		return nil, nil
	}

	if err := h.sessUseCase.DeleteSession(headerToken); err != nil {
		h.logUseCase.Report(ctx, log, msg.ErrorCannotDeleteSession(err))
		return nil, err
	}

	if timedate.IsDateExpired(session.CreationDate, h.cfg.SessionTTLHours) {
		return nil, nil
	}

	return session, nil
}
