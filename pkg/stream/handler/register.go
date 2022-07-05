package handler

import (
	"github.com/gin-gonic/gin"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	sconfig "github.com/mirumyantsev/video_hosting/pkg/config_stream"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/stream"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func RegisterTemplateHTTPEndpoints(router *gin.Engine, cfg *config.Config, scfg *sconfig.Config, uc stream.StreamUseCase,
	uuc user.UserUseCase, luc logger.LogUseCase, auc auth.AuthUseCase, suc sess.SessUseCase) {
	h := NewStreamHandler(cfg, scfg, uc, uuc, luc, auc, suc)

	router.GET("/stream", h.ServeIndex)
	router.GET("/stream/:uuid", h.ServeStream)
}

func RegisterStreamingHTTPEndpoints(router *gin.Engine, cfg *config.Config, scfg *sconfig.Config, uc stream.StreamUseCase,
	uuc user.UserUseCase, luc logger.LogUseCase, auc auth.AuthUseCase, suc sess.SessUseCase) {
	h := NewStreamHandler(cfg, scfg, uc, uuc, luc, auc, suc)

	streamRoute := router.Group("/stream")
	{
		streamRoute.GET("/codec/:uuid", h.ServeStreamCodec)
		streamRoute.POST("/receiver/:uuid", h.ServeStreamVidOverWebRTC)
		streamRoute.POST("/", h.ServeStreamWebRTC2)

		streamRoute.GET("/get/:id", h.GetStream)
		streamRoute.GET("/get/all", h.GetAllStreams)
	}
}
