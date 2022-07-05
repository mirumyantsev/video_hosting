package messages

import (
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func WarningCannotConvertCvar(cvarName string, setValue interface{}) *logger.Log {
	return &logger.Log{ErrCode: 10, Message: "Cannot convert cvar " + cvarName + ". Set default value: " + setValue.(string), ErrLevel: logger.ErrLevelWarning}
}

func FatalFailedToLoadConfigFile(err error) *logger.Log {
	return &logger.Log{ErrCode: 11, Message: "Failed to load config file. Error: " + err.Error(), ErrLevel: logger.ErrLevelFatal}
}

func InfoConfigLoaded() *logger.Log {
	return &logger.Log{Message: "Config loaded"}
}
