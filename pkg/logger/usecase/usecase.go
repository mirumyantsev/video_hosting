package usecase

import (
	"github.com/gin-gonic/gin"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/responder"
)

type LogUseCase struct {
	logRepo logger.LogRepository
}

func NewLogUseCase(logRepo logger.LogRepository) *LogUseCase {
	return &LogUseCase{
		logRepo: logRepo,
	}
}

func (u *LogUseCase) Report(ctx *gin.Context, log *logger.Log, messageLog *logger.Log) {
	logger.Complete(log, messageLog)
	logger.Finish(log)
	responder.Response(ctx, log)
	if err := u.CreateLogRecord(log); err != nil {
		logger.Complete(log, msg.ErrorCannotDoLogging(err))
		responder.Response(ctx, log)
	}
	logger.Print(log)
}

func (u *LogUseCase) ReportWithToken(ctx *gin.Context, log *logger.Log, messageLog *logger.Log, token string) {
	logger.Complete(log, messageLog)
	responder.ResponseToken(ctx, log, token)
	if err := u.CreateLogRecord(log); err != nil {
		logger.Complete(log, msg.ErrorCannotDoLogging(err))
		responder.ResponseToken(ctx, log, token)
	}
	logger.Print(log)
}

func (u *LogUseCase) CreateLogRecord(log *logger.Log) error {
	log.Message = logger.ParseMessage(log)
	return u.logRepo.CreateLogRecord(log)
}
