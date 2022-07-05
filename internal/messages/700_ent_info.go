package messages

import (
	"github.com/mirumyantsev/video_hosting/internal/info"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func ErrorStreamCannotBeEmpty() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 700, Message: "Stream cannot be empty", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCreateInfo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 701, Message: "Cannot create info. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoInfoCreated() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Info created"}
}

func ErrorCannotCheckInfoExistence(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 702, Message: "Cannot check info existence. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorInfoWithRequestedIDIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 703, Message: "Info with requested ID is not exist", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotGetInfo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 704, Message: "Cannot get info. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGotInfo(nfo *info.Info) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: nfo}
}

func ErrorCannotGetAllInfos(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 705, Message: "Cannot get all infos. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoNoInfosAvailable() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "No infos available"}
}

func InfoGotAllInfos(users map[int]*info.Info) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: users}
}

func ErrorCannotPartiallyUpdateInfo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 706, Message: "Cannot partially update info. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoInfoPartiallyUpdated() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Info partially updated"}
}

func ErrorCannotDeleteInfo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 707, Message: "Cannot delete info. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoInfoDeleted() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Info deleted"}
}
