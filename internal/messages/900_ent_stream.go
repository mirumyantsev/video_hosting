package messages

import (
	"github.com/deepch/vdk/av"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/stream"
)

func ErrorTrackIsIgnoredCodecNotSupportedWebRTC(codecType av.CodecType) *logger.Log {
	return &logger.Log{ErrCode: 900, Message: "Track is ignored - codec not supported WebRTC. codec type: " + codecType.String(), ErrLevel: logger.ErrLevelError}
}

func ErrorWritingOfCodecError(err error) *logger.Log {
	return &logger.Log{ErrCode: 901, Message: "Writing of codec error. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoStreamNotFound(suuid string) *logger.Log {
	return &logger.Log{Message: "Stream not found. Suuid: " + suuid}
}

func InfoStreamCodecNotFound(suuid string) *logger.Log {
	return &logger.Log{Message: "Stream codec not found. Suuid: " + suuid}
}

func ErrorWriteHeaderError(err error) *logger.Log {
	return &logger.Log{ErrCode: 902, Message: "WriteHeader error. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorCannotWriteBytes(err error) *logger.Log {
	return &logger.Log{ErrCode: 903, Message: "Cannot write bytes. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorStreamCodecNotFound(err error) *logger.Log {
	return &logger.Log{ErrCode: 904, Message: "Stream codec not found. LastError: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorMuxerWriteHeaderError(err error) *logger.Log {
	return &logger.Log{ErrCode: 905, Message: "Muxer WriteHeader error. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoStreamTriesToConnect(name string) *logger.Log {
	return &logger.Log{Message: "Stream tries to connect " + name}
}

func ErrorRTSPWorkerError(err error) *logger.Log {
	return &logger.Log{ErrCode: 906, Message: "rtspWorker error. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorOnDemandANDNotHasViewerError(errmsg string) *logger.Log {
	return &logger.Log{ErrCode: 907, Message: "onDemand && notHasViewer error. Error: " + errmsg, ErrLevel: logger.ErrLevelError}
}

func ErrorFrameDecoderSingleError(err error) *logger.Log {
	return &logger.Log{ErrCode: 908, Message: "frameDecoderSingle error. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoSnapshotCreated(name string) *logger.Log {
	return &logger.Log{Message: "Snapshot for " + name + " created"}
}

func ErrorBadVideoCodecWaitingForSPS_PPS() *logger.Log {
	return &logger.Log{ErrCode: 909, Message: "Bad video codec - waiting for SPS/PPS", ErrLevel: logger.ErrLevelError}
}

func InfoNoVideo() *logger.Log {
	return &logger.Log{Message: "No video"}
}

func ErrorWritePacketError(err error) *logger.Log {
	return &logger.Log{ErrCode: 910, Message: "WritePacket error. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorPseudoUUIDReadError(err error) *logger.Log {
	return &logger.Log{ErrCode: 911, Message: "pseudoUUID read error. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCheckStreamExistence(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 930, Message: "Cannot check stream existence. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorStreamWithRequestedIDIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 931, Message: "Stream with requested ID is not exist", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotGetStream(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 932, Message: "Cannot get stream. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGotStream(strm *stream.Stream) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: strm}
}

func ErrorCannotGetAllStreams(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 933, Message: "Cannot get all streams. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoNoStreamsAvailable() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "No streams available"}
}

func InfoGotAllStreams(streams map[int]*stream.Stream) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: streams}
}

func ErrorCannotGetAllWorkingStreams(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 934, Message: "Cannot get all working streams. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}
