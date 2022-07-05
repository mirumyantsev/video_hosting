package repository

import (
	"fmt"

	"github.com/mirumyantsev/video_hosting/internal/constants"
	"github.com/mirumyantsev/video_hosting/internal/info"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type InfoRepository struct {
	cfg *config.Config
}

func NewInfoRepository(cfg *config.Config) *InfoRepository {
	return &InfoRepository{cfg: cfg}
}

func (r *InfoRepository) CreateInfo(nfo *info.Info) error {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s, %s, %s, %s)", info.TableName,
		info.CreateDate, info.Stream, info.StartPeriod, info.StopPeriod,
		info.TimeLife)
	val := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s')",
		nfo.CreateDate, nfo.Stream, nfo.StartPeriod, nfo.StopPeriod,
		nfo.TimeLife)
	query := fmt.Sprintf(template, tbl, val)
	if _, err := dbo.Query(query); err != nil {
		return err
	}

	return nil
}

func (r *InfoRepository) GetInfo(id int) (*info.Info, error) {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s, %s, %s, %s", info.Id, info.Stream,
		info.StartPeriod, info.StopPeriod, info.TimeLife)
	tbl := info.TableName
	cnd := fmt.Sprintf("%s=%d", info.Id, id)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var nfo info.Info
	if err := dbo.Get(&nfo, query); err != nil {
		return nil, err
	}

	return &nfo, nil
}

func (r *InfoRepository) GetAllInfos(urlparams *user.Pagin) (map[int]*info.Info, error) {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := "*"
	tbl := info.TableName
	cnd := info.Id
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := dbo.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var infos = map[int]*info.Info{}
	var nfo info.Info
	for rows.Next() {
		if err := rows.Scan(&nfo.Id, &nfo.CreateDate, &nfo.Stream, &nfo.StartPeriod,
			&nfo.StopPeriod, &nfo.TimeLife, &nfo.UserId); err != nil {
			return nil, err
		}
		infos[nfo.Id] = &info.Info{Id: nfo.Id, CreateDate: nfo.CreateDate,
			Stream: nfo.Stream, StartPeriod: nfo.StartPeriod, StopPeriod: nfo.StopPeriod,
			TimeLife: nfo.TimeLife, UserId: nfo.UserId}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return nil, nil
	}

	return infos, nil
}

func (r *InfoRepository) PartiallyUpdateInfo(nfo *info.Info) error {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := info.TableName
	val := fmt.Sprintf("%s=CASE WHEN '%s' <> '' THEN '%s' ELSE %s END",
		info.Stream, nfo.Stream, nfo.Stream, info.Stream)
	cnd := fmt.Sprintf("%s=%d", info.Id, nfo.Id)
	query := fmt.Sprintf(template, tbl, val, cnd)

	rows, err := dbo.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *InfoRepository) DeleteInfo(id int) error {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.DELETE_FROM_TBL_WHERE_CND
	tbl := info.TableName
	cnd := fmt.Sprintf("%s=%d", info.Id, id)
	query := fmt.Sprintf(template, tbl, cnd)

	rows, err := dbo.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *InfoRepository) IsInfoExists(id int) (bool, error) {
	r.cfg.DBOName = constants.DBO_WWW_Name
	dbo := db_connect.CreateOuterDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, dbo)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := info.Id
	tbl := info.TableName
	cnd := fmt.Sprintf("%s=%d", info.Id, id)
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
