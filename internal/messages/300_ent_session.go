package messages

import (
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func ErrorCannotDeleteSession(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 300, Message: "Cannot delete session. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCreateSession(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 301, Message: "Cannot create session. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorCannotGetSessionAndDate(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 302, Message: "Cannot get session and date. Error: " + err.Error()}
}
