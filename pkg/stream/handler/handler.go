package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/deepch/vdk/av"
	webrtc "github.com/deepch/vdk/format/webrtcv3"
	"github.com/gin-gonic/gin"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	sconfig "github.com/mirumyantsev/video_hosting/pkg/config_stream"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/stream"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type StreamHandler struct {
	cfg         *config.Config
	scfg        *sconfig.Config
	useCase     stream.StreamUseCase
	userUseCase user.UserUseCase
	logUseCase  logger.LogUseCase
	authUseCase auth.AuthUseCase
	sessUseCase sess.SessUseCase
}

func NewStreamHandler(cfg *config.Config, scfg *sconfig.Config, useCase stream.StreamUseCase,
	userUseCase user.UserUseCase, logUseCase logger.LogUseCase, authUseCase auth.AuthUseCase,
	sessUseCase sess.SessUseCase) *StreamHandler {
	return &StreamHandler{
		cfg:         cfg,
		scfg:        scfg,
		useCase:     useCase,
		userUseCase: userUseCase,
		logUseCase:  logUseCase,
		authUseCase: authUseCase,
		sessUseCase: sessUseCase,
	}
}

func (h *StreamHandler) ServeIndex(ctx *gin.Context) {
	_, list := h.useCase.List()
	if len(list) > 0 {
		ctx.Header("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Redirect(http.StatusMovedPermanently, "stream/"+list[0])
	} else {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"port":    strconv.Itoa(h.cfg.ServerPort),
			"version": time.Now().String(),
		})
	}
}

func (h *StreamHandler) ServeStream(ctx *gin.Context) {
	_, list := h.useCase.List()
	sort.Strings(list)
	ctx.HTML(http.StatusOK, "player.tmpl", gin.H{
		"port":     strconv.Itoa(h.cfg.ServerPort),
		"suuid":    ctx.Param("uuid"),
		"suuidMap": list,
		"version":  time.Now().String(),
	})
}

func (h *StreamHandler) ServeStreamCodec(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if !h.useCase.Exit(uuid) {
		return
	}

	h.useCase.RunIfNotRun(uuid)

	codecs := h.useCase.CodecGet(uuid)
	if codecs == nil {
		return
	}

	var tmpCodec []stream.JCodec
	for _, codec := range codecs {
		if codec.Type() != av.H264 && codec.Type() != av.PCM_ALAW && codec.Type() != av.PCM_MULAW && codec.Type() != av.OPUS {
			logger.Printc(ctx, msg.ErrorTrackIsIgnoredCodecNotSupportedWebRTC(codec.Type()))
			continue
		}

		if codec.Type().IsVideo() {
			tmpCodec = append(tmpCodec, stream.JCodec{Type: "video"})
		} else {
			tmpCodec = append(tmpCodec, stream.JCodec{Type: "audio"})
		}
	}

	b, err := json.Marshal(tmpCodec)
	if err != nil {
		return
	}

	_, err = ctx.Writer.Write(b)
	if err != nil {
		logger.Printc(ctx, msg.ErrorWritingOfCodecError(err))
	}
}

func (h *StreamHandler) ServeStreamVidOverWebRTC(ctx *gin.Context) {
	suuid := ctx.PostForm("suuid")
	if !h.useCase.Exit(suuid) {
		logger.Printc(ctx, msg.InfoStreamNotFound(suuid))
		return
	}

	h.useCase.RunIfNotRun(suuid)

	codecs := h.useCase.CodecGet(suuid)
	if codecs == nil {
		logger.Printc(ctx, msg.InfoStreamCodecNotFound(suuid))
		return
	}

	audioOnly := false
	if len(codecs) == 1 && codecs[0].Type().IsAudio() {
		audioOnly = true
	}

	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{ICEServers: h.useCase.GetICEServers(),
		ICEUsername: h.useCase.GetICEUsername(), ICECredential: h.useCase.GetICECredential(),
		PortMin: h.useCase.GetWebRTCPortMin(), PortMax: h.useCase.GetWebRTCPortMax()})
	answer, err := muxerWebRTC.WriteHeader(codecs, ctx.PostForm("data"))
	if err != nil {
		logger.Printc(ctx, msg.ErrorWriteHeaderError(err))
		return
	}

	if _, err := ctx.Writer.Write([]byte(answer)); err != nil {
		logger.Printc(ctx, msg.ErrorCannotWriteBytes(err))
		return
	}

	go h.useCase.WritePackets(suuid, muxerWebRTC, audioOnly)
}

func (h *StreamHandler) ServeStreamWebRTC2(ctx *gin.Context) {
	url := ctx.PostForm("url")
	if _, ok := h.scfg.Streams[url]; !ok {
		h.scfg.Streams[url] = sconfig.Stream{
			URL:        url,
			OnDemand:   true,
			ClientList: make(map[string]sconfig.Viewer),
		}
	}

	h.useCase.RunIfNotRun(url)

	codecs := h.useCase.CodecGet(url)
	if codecs == nil {
		logger.Printc(ctx, msg.ErrorStreamCodecNotFound(h.scfg.LastError))
		return
	}

	muxerWebRTC := webrtc.NewMuxer(
		webrtc.Options{
			ICEServers: h.useCase.GetICEServers(),
			PortMin:    h.useCase.GetWebRTCPortMin(),
			PortMax:    h.useCase.GetWebRTCPortMax(),
		},
	)

	sdp64 := ctx.PostForm("sdp64")
	answer, err := muxerWebRTC.WriteHeader(codecs, sdp64)
	if err != nil {
		logger.Printc(ctx, msg.ErrorMuxerWriteHeaderError(err))
		return
	}

	response := stream.Response{
		Sdp64: answer,
	}

	for _, codec := range codecs {
		if codec.Type() != av.H264 &&
			codec.Type() != av.PCM_ALAW &&
			codec.Type() != av.PCM_MULAW &&
			codec.Type() != av.OPUS {
			logger.Printc(ctx, msg.ErrorTrackIsIgnoredCodecNotSupportedWebRTC(codec.Type()))
			continue
		}
		if codec.Type().IsVideo() {
			response.Tracks = append(response.Tracks, "video")
		} else {
			response.Tracks = append(response.Tracks, "audio")
		}
	}

	ctx.JSON(200, response)

	audioOnly := len(codecs) == 1 && codecs[0].Type().IsAudio()

	go h.useCase.WritePackets(url, muxerWebRTC, audioOnly)
}
