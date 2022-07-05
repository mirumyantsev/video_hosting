package repository

import (
	"fmt"

	"github.com/mirumyantsev/video_hosting/internal/constants"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/stream"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type StreamRepository struct {
	cfg *config.Config
}

func NewStreamRepository(cfg *config.Config) *StreamRepository {
	return &StreamRepository{cfg: cfg}
}

func (r *StreamRepository) GetStream(id int) (*stream.StreamGet, error) {
	r.cfg.DBOName = constants.DBO_L3_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", stream.Id,
		stream.StreamColumn, stream.DateTime, stream.StatePublic,
		stream.StatusPublic, stream.StatusRecord, stream.PathStream)
	tbl := stream.TableName
	cnd := fmt.Sprintf("%s=%d", stream.Id, id)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var strm stream.StreamGet
	if err := dbo.Get(&strm, query); err != nil {
		return nil, err
	}

	return &strm, nil
}

func (r *StreamRepository) GetAllStreams(urlparams *user.Pagin) (map[int]*stream.StreamGet, error) {
	r.cfg.DBOName = constants.DBO_L3_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := fmt.Sprintf("%s, %s, %s, %s, %s", stream.Id,
		stream.StreamColumn, stream.DateTime,
		stream.StatusPublic, stream.PathStream)
	tbl := stream.TableName
	cnd := stream.Id
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := dbo.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streams = map[int]*stream.StreamGet{}
	var strm stream.StreamGet
	for rows.Next() {
		if err := rows.Scan(&strm.Id, &strm.Stream, &strm.DateTime,
			&strm.StatusPublic, &strm.PathStream); err != nil {
			return nil, err
		}
		streams[strm.Id] = &stream.StreamGet{Id: strm.Id, Stream: strm.Stream,
			DateTime: strm.DateTime, StatusPublic: strm.StatusPublic,
			PathStream: strm.PathStream}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(streams) == 0 {
		return nil, nil
	}

	return streams, nil
}

func (r *StreamRepository) GetAllWorkingStreams() (*[]string, error) {
	r.cfg.DBOName = constants.DBO_L3_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := stream.StreamColumn
	tbl := stream.TableName
	cnd := fmt.Sprintf("%s=1", stream.StatusPublic)
	query := fmt.Sprintf(template, col, tbl, cnd)

	rows, err := dbo.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workingStreams := []string{}
	strm := ""
	for rows.Next() {
		if err := rows.Scan(&strm); err != nil {
			return nil, err
		}
		workingStreams = append(workingStreams, strm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(workingStreams) == 0 {
		return nil, nil
	}

	return &workingStreams, nil
}

func (r *StreamRepository) IsStreamExists(id int) (bool, error) {
	r.cfg.DBOName = constants.DBO_L3_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := stream.Id
	tbl := stream.TableName
	cnd := fmt.Sprintf("%s=%d", stream.Id, id)
	query := fmt.Sprintf(template, col, tbl, cnd)
	rows, err := dbo.Query(query)

	if err != nil {
		return false, err
	}
	defer rows.Close()

	if isRowPresent := rows.Next(); !isRowPresent {
		return false, nil
	}

	return true, nil
}
