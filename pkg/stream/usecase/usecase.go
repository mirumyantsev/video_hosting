package usecase

import (
	"errors"
	"fmt"

	"image/jpeg"
	"os"
	"time"

	"github.com/deepch/vdk/av"

	"github.com/deepch/vdk/cgo/ffmpeg"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtspv2"
	webrtc "github.com/deepch/vdk/format/webrtcv3"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	sconfig "github.com/mirumyantsev/video_hosting/pkg/config_stream"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/stream"
)

const (
	errorStreamExitNoViewer = "Stream was exited on demand - no Viewer"
	snapshotPath            = "./media/%s/images"
	snapshotName            = "snapshot.jpg"
	videoTimeoutSeconds     = 80
)

type StreamUseCase struct {
	cfg        *config.Config
	scfg       *sconfig.Config
	streamRepo stream.StreamRepository
}

func NewStreamUseCase(cfg *config.Config, scfg *sconfig.Config, streamRepo stream.StreamRepository) *StreamUseCase {
	return &StreamUseCase{
		cfg:        cfg,
		scfg:       scfg,
		streamRepo: streamRepo,
	}
}

func (u *StreamUseCase) ServeStreams() {
	u.scfg.Streams = map[string]sconfig.Stream{}

	go func() {
		for {
			u.startDefaultGetWorkingStreamsCycle()
			for {
				time.Sleep(1 * time.Second)
				if u.scfg.StreamDropped {
					u.scfg.StreamDropped = false
					break
				}
			}
		}
	}()
}

func (u *StreamUseCase) startDefaultGetWorkingStreamsCycle() {
	for {
		workingStreams, err := u.getAllWorkingStreams()
		if err != nil {
			logger.Printc(nil, msg.ErrorCannotGetAllWorkingStreams(err))
			return
		}
		for _, url := range *workingStreams {
			if _, found := u.scfg.Streams[url]; !found {
				u.scfg.Streams[url] = sconfig.Stream{ClientList: make(map[string]sconfig.Viewer), URL: url}
				if !u.scfg.Streams[url].OnDemand {
					go u.rtspWorkerLoop(url, u.scfg.Streams[url].URL, u.scfg.Streams[url].OnDemand, u.scfg.Streams[url].DisableAudio, u.scfg.Streams[url].Debug)
					u.scfg.StreamsCount++
				}
			}
		}
		go func() {
			time.Sleep(200 * time.Millisecond)
			logger.Printc(nil, &logger.Log{Message: fmt.Sprintf("Working streams: %d", u.scfg.StreamsCount)})
		}()
		time.Sleep(time.Duration(u.cfg.StreamStreamsUpdatePeriodSeconds) * time.Second)
	}
}

func (u *StreamUseCase) rtspWorkerLoop(name, url string, onDemand, disableAudio, debug bool) {
	defer func() {
		if _, found := u.scfg.Streams[name]; found {
			delete(u.scfg.Streams, name)
			u.scfg.StreamDropped = true
			u.scfg.StreamsCount--
			logger.Printc(nil, &logger.Log{Message: "Stream dropped. Stream: " + name})
		}
	}()
	defer u.runUnlock(name)
	for {
		logger.Printc(nil, msg.InfoStreamTriesToConnect(name))
		err := u.rtspWorker(name, url, onDemand, disableAudio, debug)
		if err != nil {
			u.scfg.LastError = err
			logger.Printc(nil, msg.ErrorRTSPWorkerError(err))
			return
		}
		if onDemand && !u.isHasViewer(name) {
			logger.Printc(nil, msg.ErrorOnDemandANDNotHasViewerError(errorStreamExitNoViewer))
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (u *StreamUseCase) rtspWorker(name, url string, onDemand, disableAudio, debug bool) error {
	keyTest := time.NewTimer(20 * time.Second)
	clientTest := time.NewTimer(20 * time.Second)

	// add next timeout
	newRTSPClient := rtspv2.RTSPClientOptions{
		URL:              os.Getenv("RTSP_URL_MAIN") + url,
		DisableAudio:     disableAudio,
		DialTimeout:      3 * time.Second,
		ReadWriteTimeout: 3 * time.Second,
		Debug:            debug,
	}

	rtspClient, err := rtspv2.Dial(newRTSPClient)
	if err != nil {
		return err
	}
	defer rtspClient.Close()

	if rtspClient.CodecData != nil {
		u.codecAdd(name, rtspClient.CodecData)
	}

	audioOnly := false
	videoIDX := 0
	for i, codec := range rtspClient.CodecData {
		if codec.Type().IsVideo() {
			audioOnly = false
			videoIDX = i
		}
	}

	var frameDecoderSingle *ffmpeg.VideoDecoder
	if !audioOnly {
		frameDecoderSingle, err = ffmpeg.NewVideoDecoder(rtspClient.CodecData[videoIDX].(av.VideoCodecData))
		if err != nil {
			logger.Printc(nil, msg.ErrorFrameDecoderSingleError(err))
		}
	}

	isTimeToSnapshot := true
	if u.cfg.StreamSnapshotsEnable {
		go func() {
			for {
				time.Sleep(time.Duration(u.cfg.StreamSnapshotPeriodSeconds) * time.Second)
				isTimeToSnapshot = true
			}
		}()
	}

	snapshotDir := fmt.Sprintf(snapshotPath, name)
	if !IsPathExists(snapshotDir) {
		os.MkdirAll(snapshotDir, 0777)
	}

	for {
		select {
		case <-clientTest.C:
			if onDemand {
				if !u.isHasViewer(name) {
					return errors.New(errorStreamExitNoViewer)
				} else {
					clientTest.Reset(20 * time.Second)
				}
			}
		case <-keyTest.C:
			return errors.New("stream exit - no video on stream")
		case signals := <-rtspClient.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				u.codecAdd(name, rtspClient.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return errors.New("stream exit - rtsp disconnect")
			}
		case packetAV := <-rtspClient.OutgoingPacketQueue:
			if audioOnly || packetAV.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
			}
			u.cast(name, *packetAV)
			// sample single frame decode encode to jpeg, save on disk
			if !u.cfg.StreamSnapshotsEnable || !packetAV.IsKeyFrame {
				break
			}
			pic, err := frameDecoderSingle.DecodeSingle(packetAV.Data)
			if err != nil ||
				pic == nil || !isTimeToSnapshot {
				break
			}
			out, err := os.Create(snapshotDir + "/" + snapshotName)
			if err != nil {
				break
			}
			if err := jpeg.Encode(out, &pic.Image, nil); err == nil {
				if u.cfg.StreamSnapshotShowStatus {
					logger.Printc(nil, msg.InfoSnapshotCreated(name))
				}
				isTimeToSnapshot = false
			}
		}
	}
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func (u *StreamUseCase) codecAdd(suuid string, codecs []av.CodecData) {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	t := u.scfg.Streams[suuid]
	t.Codecs = codecs
	u.scfg.Streams[suuid] = t
}

func (u *StreamUseCase) isHasViewer(uuid string) bool {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	if cfg, ok := u.scfg.Streams[uuid]; ok && len(cfg.ClientList) > 0 {
		return true
	}
	return false
}

func (u *StreamUseCase) cast(uuid string, pck av.Packet) {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	for _, val := range u.scfg.Streams[uuid].ClientList {
		if len(val.Cast) < cap(val.Cast) {
			val.Cast <- pck
		}
	}
}

func (u *StreamUseCase) runUnlock(uuid string) {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	cfg, ok := u.scfg.Streams[uuid]
	if !ok {
		return
	}
	if cfg.OnDemand && cfg.RunLock {
		cfg.RunLock = false
		u.scfg.Streams[uuid] = cfg
	}
}

func (u *StreamUseCase) Exit(suuid string) bool {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	_, ok := u.scfg.Streams[suuid]
	return ok
}

func (u *StreamUseCase) RunIfNotRun(uuid string) {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	cfg, ok := u.scfg.Streams[uuid]
	if !ok {
		return
	}
	if cfg.OnDemand && !cfg.RunLock {
		cfg.RunLock = true
		u.scfg.Streams[uuid] = cfg
		go u.rtspWorkerLoop(uuid, cfg.URL, cfg.OnDemand, cfg.DisableAudio, cfg.Debug)
	}
}

func (u *StreamUseCase) CodecGet(suuid string) []av.CodecData {
	for i := 0; i < 100; i++ {
		u.scfg.StreamsMutex.RLock()
		cfg, ok := u.scfg.Streams[suuid]
		u.scfg.StreamsMutex.RUnlock()
		if !ok {
			return nil
		}
		if cfg.Codecs == nil {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		for _, codec := range cfg.Codecs {
			if codec.Type() != av.H264 {
				continue
			}
			codecVideo := codec.(h264parser.CodecData)
			if codecVideo.SPS() == nil && codecVideo.PPS() == nil &&
				len(codecVideo.SPS()) <= 0 && len(codecVideo.PPS()) <= 0 {
				logger.Printc(nil, msg.ErrorBadVideoCodecWaitingForSPS_PPS())
				time.Sleep(50 * time.Millisecond)
			}
		}
		return cfg.Codecs
	}
	return nil
}

func (u *StreamUseCase) GetICEServers() []string {
	u.cfg.StreamICEServersMutex.Lock()
	defer u.cfg.StreamICEServersMutex.Unlock()
	return u.cfg.StreamICEServers
}

func (u *StreamUseCase) GetICEUsername() string {
	u.scfg.Server.ICEUsernameMutex.Lock()
	defer u.scfg.Server.ICEUsernameMutex.Unlock()
	return u.scfg.Server.ICEUsername
}

func (u *StreamUseCase) GetICECredential() string {
	u.scfg.Server.ICECredentialMutex.Lock()
	defer u.scfg.Server.ICECredentialMutex.Unlock()
	return u.scfg.Server.ICECredential
}

func (u *StreamUseCase) GetWebRTCPortMin() uint16 {
	u.scfg.Server.WebRTCPortMinMutex.Lock()
	defer u.scfg.Server.WebRTCPortMinMutex.Unlock()
	return u.scfg.Server.WebRTCPortMin
}

func (u *StreamUseCase) GetWebRTCPortMax() uint16 {
	u.scfg.Server.WebRTCPortMaxMutex.Lock()
	defer u.scfg.Server.WebRTCPortMaxMutex.Unlock()
	return u.scfg.Server.WebRTCPortMax
}

func (u *StreamUseCase) WritePackets(url string, muxerWebRTC *webrtc.Muxer, audioOnly bool) {
	cid, ch := u.CastListAdd(url)
	defer u.CastListDelete(url, cid)
	defer muxerWebRTC.Close()
	videoStart := false
	noVideo := time.NewTimer(videoTimeoutSeconds * time.Second)
	for {
		select {
		case <-noVideo.C:
			logger.Printc(nil, msg.InfoNoVideo())
			return
		case pck := <-ch:
			if pck.IsKeyFrame || audioOnly {
				noVideo.Reset(videoTimeoutSeconds * time.Second)
				videoStart = true
			}
			if !videoStart && !audioOnly {
				continue
			}
			err := muxerWebRTC.WritePacket(pck)
			if err != nil {
				logger.Printc(nil, msg.ErrorWritePacketError(err))
				return
			}
		}
	}
}

func (u *StreamUseCase) CastListAdd(suuid string) (string, chan av.Packet) {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	ch := make(chan av.Packet, 100)
	u.scfg.Streams[suuid].ClientList[suuid] = sconfig.Viewer{Cast: ch}
	return suuid, ch
}

func (u *StreamUseCase) CastListDelete(suuid, cuuid string) {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	delete(u.scfg.Streams[suuid].ClientList, cuuid)
}

func (u *StreamUseCase) List() (string, []string) {
	u.scfg.StreamsMutex.Lock()
	defer u.scfg.StreamsMutex.Unlock()
	var res []string
	var first string
	for key := range u.scfg.Streams {
		if first == "" {
			first = key
		}
		res = append(res, key)
	}
	return first, res
}
