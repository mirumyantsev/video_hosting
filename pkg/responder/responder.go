package responder

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

type MessageOutput struct {
	Message string `json:"message"`
}

type MessageTokenOutput struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ErrorOutput struct {
	Error interface{} `json:"error"`
}

type ErrorData struct {
	ErrCode   int    `json:"errCode"`
	Statement string `json:"statement"`
}

func ResponseToken(ctx *gin.Context, log *logger.Log, token string) {
	ctx.AbortWithStatusJSON(log.StatusCode, MessageTokenOutput{Message: log.Message.(string), Token: token})
}

func Response(ctx *gin.Context, log *logger.Log) {
	if reflect.TypeOf(log.Message) == reflect.TypeOf("") {
		if log.ErrLevel != logger.ErrLevelInfo {
			ctx.AbortWithStatusJSON(log.StatusCode, ErrorOutput{
				ErrorData{ErrCode: log.ErrCode, Statement: log.Message.(string)},
			})
			return
		}

		ctx.AbortWithStatusJSON(log.StatusCode, MessageOutput{log.Message.(string)})
		return
	}

	ctx.AbortWithStatusJSON(log.StatusCode, log.Message)
}
