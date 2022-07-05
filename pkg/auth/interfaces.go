package auth

import (
	"github.com/gin-gonic/gin"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
)

type AuthCommon interface {
	GetNamepass(namepass *Namepass) error
	UpdateUserPassword(namepass *Namepass) error
	IsUsernameAndPasswordExists(username, passwordHash string) (bool, error)
}

type AuthUseCase interface {
	AuthCommon

	ReadHeader(ctx *gin.Context) string
	IsTokenExists(token string) bool
	IsMatched(username_1, username_2 string) bool
	IsRequiredEmpty(namepass *Namepass) bool
	IsSessionExists(session *sess.Session) bool
	BindJSONNamepass(ctx *gin.Context) (*Namepass, error)
	GenerateToken(namepass *Namepass) (string, error)
	ParseToken(token string) (*Namepass, error)
}

type AuthRepository interface {
	AuthCommon

	UpdateNamepassLastLogin(username, token string) error
}
