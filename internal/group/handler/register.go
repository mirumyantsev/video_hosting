package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mikerumy/vhosting/internal/group"
	lg "github.com/mikerumy/vhosting/internal/logging"
	sess "github.com/mikerumy/vhosting/internal/session"
	"github.com/mikerumy/vhosting/pkg/auth"
	"github.com/mikerumy/vhosting/pkg/config"
	"github.com/mikerumy/vhosting/pkg/user"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, uc group.GroupUseCase, luc lg.LogUseCase,
	auc auth.AuthUseCase, suc sess.SessUseCase, uuc user.UserUseCase) {
	h := NewGroupHandler(cfg, uc, luc, auc, suc, uuc)

	groupRoute := router.Group("/group")
	{
		groupRoute.POST("", h.CreateGroup)
		groupRoute.GET(":id", h.GetGroup)
		groupRoute.GET("all", h.GetAllGroups)
		groupRoute.PATCH(":id", h.PartiallyUpdateGroup)
		groupRoute.DELETE(":id", h.DeleteGroup)
	}

	groupSetUserRoute := router.Group("/group/user")
	{
		groupSetUserRoute.POST(":id", h.SetUserGroups)
		groupSetUserRoute.GET(":id", h.GetUserGroups)
		groupSetUserRoute.DELETE(":id", h.DeleteUserGroups)
	}
}
