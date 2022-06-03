package handler

import (
	"github.com/gin-gonic/gin"
	lg "github.com/mikerumy/vhosting/internal/logging"
	sess "github.com/mikerumy/vhosting/internal/session"
	"github.com/mikerumy/vhosting/pkg/auth"
	"github.com/mikerumy/vhosting/pkg/config"
	"github.com/mikerumy/vhosting/pkg/user"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, uc auth.AuthUseCase, uuc user.UserUseCase,
	suc sess.SessUseCase, luc lg.LogUseCase) {
	h := NewAuthHandler(cfg, uc, uuc, suc, luc)

	authRoute := router.Group("/auth")
	{
		authRoute.POST("/signin", h.SignIn)
		authRoute.POST("/change_password", h.ChangePassword)
		authRoute.GET("/signout", h.SignOut)
	}
}
