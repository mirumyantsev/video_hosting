package repository

import (
	"fmt"

	perm "github.com/mirumyantsev/video_hosting/internal/permission"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func (r *PermRepository) SetGroupPermissions(values string) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s)", perm.GPTableName, perm.GroupId,
		perm.PermId)
	val := values
	query := fmt.Sprintf(template, tbl, val)

	if _, err := db.Query(query); err != nil {
		return err
	}

	return nil
}

func (r *PermRepository) GetGroupPermissions(id int, urlparams *user.Pagin) (*perm.PermIds, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := perm.PermId
	tbl := perm.GPTableName
	cnd := fmt.Sprintf("%s=$1 AND %s", perm.GroupId, perm.Id)
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permIds perm.PermIds
	var num int
	for rows.Next() {
		if err := rows.Scan(&num); err != nil {
			return nil, err
		}
		permIds.Ids = append(permIds.Ids, num)
	}

	return &permIds, nil
}

func (r *PermRepository) DeleteGroupPermissions(id int, condIds string) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.DELETE_FROM_TBL_WHERE_CND
	tbl := perm.GPTableName
	cnd := fmt.Sprintf("%s=$1 AND %s IN (%s)", perm.GroupId, perm.PermId,
		condIds)
	query := fmt.Sprintf(template, tbl, cnd)

	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
