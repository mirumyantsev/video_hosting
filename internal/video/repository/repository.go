package repository

import (
	"fmt"

	"github.com/mirumyantsev/video_hosting/internal/constants"
	"github.com/mirumyantsev/video_hosting/internal/video"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type VideoRepository struct {
	cfg *config.Config
}

func NewVideoRepository(cfg *config.Config) *VideoRepository {
	return &VideoRepository{cfg: cfg}
}

func (r *VideoRepository) CreateVideo(vid *video.Video) error {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s, %s, %s, %s)", video.TableName,
		video.Url, video.File, video.CreateDate, video.InfoId,
		video.UserId)
	val := fmt.Sprintf("('%s', '%s', '%s', %d, %d)", vid.Url, vid.File,
		vid.CreateDate, vid.InfoId, vid.UserId)
	query := fmt.Sprintf(template, tbl, val)

	if _, err := dbo.Query(query); err != nil {
		return err
	}

	return nil
}

func (r *VideoRepository) GetVideo(id int) (*video.Video, error) {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s, %s, %s, %s, %s", video.Id, video.Url,
		video.File, video.CreateDate, video.InfoId, video.UserId)
	tbl := video.TableName
	cnd := fmt.Sprintf("%s=%d", video.Id, id)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var vid video.Video
	if err := dbo.Get(&vid, query); err != nil {
		return nil, err
	}

	return &vid, nil
}

func (r *VideoRepository) GetAllVideos(urlparams *user.Pagin) (map[int]*video.Video, error) {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := "*"
	tbl := video.TableName
	cnd := video.Id
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := dbo.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos = map[int]*video.Video{}
	var vid video.Video
	for rows.Next() {
		if err := rows.Scan(&vid.Id, &vid.Url, &vid.File, &vid.CreateDate,
			&vid.InfoId, &vid.UserId); err != nil {
			return nil, err
		}
		videos[vid.Id] = &video.Video{Id: vid.Id, Url: vid.Url,
			File: vid.File, CreateDate: vid.CreateDate,
			InfoId: vid.InfoId, UserId: vid.UserId}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(videos) == 0 {
		return nil, nil
	}

	return videos, nil
}

func (r *VideoRepository) PartiallyUpdateVideo(vid *video.Video) error {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := video.TableName
	val := fmt.Sprintf("%s=CASE WHEN '%s' <> '' THEN '%s' ELSE %s END, ", video.Url, vid.Url, vid.Url, video.Url) +
		fmt.Sprintf("%s=CASE WHEN '%s' <> '' THEN '%s' ELSE %s END, ", video.File, vid.File, vid.File, video.File) +
		fmt.Sprintf("%s=CASE WHEN %d > -1 THEN %d ELSE %s END", video.InfoId, vid.InfoId, vid.InfoId, video.InfoId)
	cnd := fmt.Sprintf("%s=%d", video.Id, vid.Id)
	query := fmt.Sprintf(template, tbl, val, cnd)

	rows, err := dbo.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *VideoRepository) DeleteVideo(id int) error {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.DELETE_FROM_TBL_WHERE_CND
	tbl := video.TableName
	cnd := fmt.Sprintf("%s=%d", video.Id, id)
	query := fmt.Sprintf(template, tbl, cnd)

	rows, err := dbo.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *VideoRepository) IsVideoExists(id int) (bool, error) {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := video.Id
	tbl := video.TableName
	cnd := fmt.Sprintf("%s=%d", video.Id, id)
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
