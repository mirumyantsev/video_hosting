package handler

import (
	"github.com/gin-gonic/gin"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, uc user.UserUseCase,
	luc logger.LogUseCase, auc auth.AuthUseCase, suc sess.SessUseCase) {
	h := NewUserHandler(cfg, uc, luc, auc, suc)

	userRoute := router.Group("/user")
	{
		userRoute.POST("", h.CreateUser)
		userRoute.GET(":id", h.GetUser)
		userRoute.GET("all", h.GetAllUsers)
		userRoute.POST("/change_password", h.UpdateUserPassword)
		userRoute.PATCH(":id", h.PartiallyUpdateUser)
		userRoute.DELETE(":id", h.DeleteUser)
	}
}
