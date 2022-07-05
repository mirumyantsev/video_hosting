package stream

import (
	"github.com/deepch/vdk/av"
	webrtc "github.com/deepch/vdk/format/webrtcv3"
	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type StreamCommon interface {
	IsStreamExists(id int) (bool, error)
}

type StreamUseCase interface {
	StreamCommon

	GetStream(id int) (*Stream, error)
	GetAllStreams(urlparams *user.Pagin) (map[int]*Stream, error)

	AtoiRequestedId(ctx *gin.Context) (int, error)
	ParseURLParams(ctx *gin.Context) *user.Pagin

	ServeStreams()
	Exit(suuid string) bool
	RunIfNotRun(uuid string)
	CodecGet(suuid string) []av.CodecData
	GetICEServers() []string
	GetICEUsername() string
	GetICECredential() string
	GetWebRTCPortMin() uint16
	GetWebRTCPortMax() uint16
	WritePackets(url string, muxerWebRTC *webrtc.Muxer, audioOnly bool)
	CastListAdd(suuid string) (string, chan av.Packet)
	CastListDelete(suuid, cuuid string)
	List() (string, []string)
}

type StreamRepository interface {
	StreamCommon

	GetStream(id int) (*StreamGet, error)
	GetAllStreams(urlparams *user.Pagin) (map[int]*StreamGet, error)
	GetAllWorkingStreams() (*[]string, error)
}
