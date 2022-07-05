package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/info"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	sconfig "github.com/mirumyantsev/video_hosting/pkg/config_stream"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, scfg *sconfig.Config, uc info.InfoUseCase, luc logger.LogUseCase,
	auc auth.AuthUseCase, suc sess.SessUseCase, uuc user.UserUseCase) {
	h := NewInfoHandler(cfg, scfg, uc, luc, auc, suc, uuc)

	infoRoute := router.Group("/info")
	{
		infoRoute.POST("", h.CreateInfo)
		infoRoute.GET(":id", h.GetInfo)
		infoRoute.GET("all", h.GetAllInfos)
		infoRoute.PATCH(":id", h.PartiallyUpdateInfo)
		infoRoute.DELETE(":id", h.DeleteInfo)

		infoRoute.GET("/test", h.Temp)
	}
}
