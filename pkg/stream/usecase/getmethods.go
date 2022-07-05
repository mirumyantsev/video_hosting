package usecase

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/pkg/stream"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func (u *StreamUseCase) GetStream(id int) (*stream.Stream, error) {
	strmg, err := u.streamRepo.GetStream(id)
	if err != nil {
		return nil, err
	}
	var strm stream.Stream
	strm.Id = strmg.Id
	strm.Stream = strmg.Stream.String
	strm.DateTime = strmg.DateTime.String
	strm.StatePublic = int(strmg.StatePublic.Int16)
	strm.StatusPublic = int(strmg.StatusPublic.Int16)
	strm.StatusRecord = int(strmg.StatusRecord.Int16)
	strm.PathStream = strmg.PathStream.String
	if (strm.Stream != "") && (strm.StatusPublic != 0) {
		u.cfg.StreamLink = os.Getenv("RTSP_URL_MAIN") + strm.Stream
		fmt.Println("Got RTSP link")
	} else {
		u.cfg.StreamLink = ""
		fmt.Println("No RTSP link")
	}
	return &strm, nil
}

func (u *StreamUseCase) GetAllStreams(urlparams *user.Pagin) (map[int]*stream.Stream, error) {
	streamsg, err := u.streamRepo.GetAllStreams(urlparams)
	if err != nil {
		return nil, err
	}
	var streams = map[int]*stream.Stream{}
	for _, val := range streamsg {
		streams[val.Id] = &stream.Stream{Id: val.Id, Stream: val.Stream.String,
			DateTime: val.DateTime.String, StatusPublic: int(val.StatusPublic.Int16),
			PathStream: val.PathStream.String}
	}
	return streams, nil
}

func (u *StreamUseCase) getAllWorkingStreams() (*[]string, error) {
	return u.streamRepo.GetAllWorkingStreams()
}

func (u *StreamUseCase) AtoiRequestedId(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (u *StreamUseCase) IsStreamExists(id int) (bool, error) {
	exists, err := u.streamRepo.IsStreamExists(id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *StreamUseCase) ParseURLParams(ctx *gin.Context) *user.Pagin {
	urlparams := ctx.Request.URL.Query()
	var pagin user.Pagin
	if lim := urlparams.Get("_limit"); lim != "" {
		pagin.Limit, _ = strconv.Atoi(lim)
	}
	if pg := urlparams.Get("_page"); pg != "" {
		pagin.Page, _ = strconv.Atoi(pg)
	}
	pagin.Page = pagin.Page*pagin.Limit - pagin.Limit
	if pagin.Limit == 0 {
		pagin.Limit = u.cfg.PaginationGetLimitDefault
	}
	return &pagin
}
