package repository

import (
	"fmt"

	"github.com/mirumyantsev/video_hosting/internal/group"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func (r *GroupRepository) SetUserGroups(values string) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s)", group.UGTableName, group.UserId,
		group.GroupId)
	val := values
	query := fmt.Sprintf(template, tbl, val)

	if _, err := db.Query(query); err != nil {
		return err
	}

	return nil
}

func (r *GroupRepository) GetUserGroups(id int, urlparams *user.Pagin) (*group.GroupIds, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := group.GroupId
	tbl := group.UGTableName
	cnd := fmt.Sprintf("%s=$1 AND %s", group.UserId, group.Id)
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupIds group.GroupIds
	var grp int
	for rows.Next() {
		if err := rows.Scan(&grp); err != nil {
			return nil, err
		}
		groupIds.Ids = append(groupIds.Ids, grp)
	}

	return &groupIds, nil
}

func (r *GroupRepository) DeleteUserGroups(id int, condIds string) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.DELETE_FROM_TBL_WHERE_CND
	tbl := group.UGTableName
	cnd := fmt.Sprintf("%s=$1 AND %s IN (%s)", group.UserId, group.GroupId,
		condIds)
	query := fmt.Sprintf(template, tbl, cnd)

	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
