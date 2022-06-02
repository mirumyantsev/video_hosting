package messages

import (
	lg "github.com/mikerumy/vhosting/internal/logging"
	"github.com/mikerumy/vhosting/pkg/download"
	"github.com/mikerumy/vhosting/pkg/logger"
)

func ErrorExtensionIsNotMp4() *lg.Log {
	return &lg.Log{StatusCode: 405, ErrCode: 1000, Message: "Extension is not .mp4.", ErrLevel: logger.ErrLevelError}
}

func InfoPutDownloadLink(dload *download.Download) *lg.Log {
	return &lg.Log{StatusCode: 200, Message: dload, ErrLevel: logger.ErrLevelInfo}
}