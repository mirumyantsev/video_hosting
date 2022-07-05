package usecase

import (
	"github.com/gin-gonic/gin"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/hasher"
	"github.com/mirumyantsev/video_hosting/pkg/headers"
)

type AuthUseCase struct {
	cfg      *config.Config
	authRepo auth.AuthRepository
}

func NewAuthUseCase(cfg *config.Config, authRepo auth.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		cfg:      cfg,
		authRepo: authRepo,
	}
}

func (u *AuthUseCase) IsTokenExists(token string) bool {
	if token != "" {
		return true
	}
	return false
}

func (u *AuthUseCase) IsRequiredEmpty(namepass *auth.Namepass) bool {
	if namepass.PasswordHash == "" {
		return true
	}
	return false
}

func (u *AuthUseCase) IsSessionExists(session *sess.Session) bool {
	if session.Content != "" {
		return true
	}
	return false
}

func (u *AuthUseCase) IsMatched(username_1, username_2 string) bool {
	if username_1 == username_2 {
		return true
	}
	return false
}

func (u *AuthUseCase) GetNamepass(namepass *auth.Namepass) error {
	return u.authRepo.GetNamepass(namepass)
}

func (u *AuthUseCase) UpdateUserPassword(namepass *auth.Namepass) error {
	return u.authRepo.UpdateUserPassword(namepass)
}

func (u *AuthUseCase) IsUsernameAndPasswordExists(username, passwordHash string) (bool, error) {
	exists, err := u.authRepo.IsUsernameAndPasswordExists(username, passwordHash)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *AuthUseCase) ReadHeader(ctx *gin.Context) string {
	return headers.ReadHeader(ctx)
}

func (u *AuthUseCase) BindJSONNamepass(ctx *gin.Context) (*auth.Namepass, error) {
	var namepass auth.Namepass
	if err := ctx.BindJSON(&namepass); err != nil {
		return &namepass, err
	}
	if namepass.PasswordHash != "" {
		namepass.PasswordHash = hasher.GeneratePasswordHash(namepass.PasswordHash, u.cfg.HashingPasswordSalt)
	}
	return &namepass, nil
}

func (u *AuthUseCase) GenerateToken(namepass *auth.Namepass) (string, error) {
	namepass.PasswordHash = hasher.GeneratePasswordHash(namepass.PasswordHash, u.cfg.HashingPasswordSalt)
	token, err := hasher.GenerateToken(namepass, u.cfg.HashingTokenSigningKey, u.cfg.SessionTTLHours)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *AuthUseCase) ParseToken(token string) (*auth.Namepass, error) {
	return hasher.ParseToken(token, u.cfg.HashingTokenSigningKey)
}
