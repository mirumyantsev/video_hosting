package messages

import (
	"github.com/mirumyantsev/video_hosting/internal/video"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func ErrorUrlAndFilenameCannotBeEmpty() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 800, Message: "URL and file name cannot be empty", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCreateVideo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 801, Message: "Cannot create video. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoVideoCreated() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Video created"}
}

func ErrorCannotCheckVideoExistence(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 802, Message: "Cannot check video existence. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorVideoWithRequestedIDIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 803, Message: "Video with requested ID is not exist", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotGetVideo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 804, Message: "Cannot get video. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGotVideo(nfo *video.Video) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: nfo}
}

func ErrorCannotGetAllVideos(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 805, Message: "Cannot get all videos. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoNoVideosAvailable() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "No videos available"}
}

func InfoGotAllVideos(users map[int]*video.Video) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: users}
}

func ErrorCannotPartiallyUpdateVideo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 806, Message: "Cannot partially update video. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoVideoPartiallyUpdated() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Video partially updated"}
}

func ErrorCannotDeleteVideo(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 807, Message: "Cannot delete video. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoVideoDeleted() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Video deleted"}
}
