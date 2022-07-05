package repository

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/mirumyantsev/video_hosting/internal/group"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type GroupRepository struct {
	cfg *config.Config
}

func NewGroupRepository(cfg *config.Config) *GroupRepository {
	return &GroupRepository{cfg: cfg}
}

func (r *GroupRepository) CreateGroup(grp *group.Group) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s)", group.TableName, group.Name)
	val := "($1)"
	query := fmt.Sprintf(template, tbl, val)

	if _, err := db.Query(query, grp.Name); err != nil {
		return err
	}

	return nil
}

func (r *GroupRepository) GetGroup(id int) (*group.Group, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s", group.Id, group.Name)
	tbl := group.TableName
	cnd := fmt.Sprintf("%s=$1", group.Id)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var grp group.Group
	if err := db.Get(&grp, query, id); err != nil {
		return nil, err
	}

	return &grp, nil
}

func (r *GroupRepository) GetAllGroups(urlparams *user.Pagin) (map[int]*group.Group, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := "*"
	tbl := group.TableName
	cnd := group.Id
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups = map[int]*group.Group{}
	var grp group.Group
	for rows.Next() {
		if err := rows.Scan(&grp.Id, &grp.Name); err != nil {
			return nil, err
		}
		groups[grp.Id] = &group.Group{Id: grp.Id, Name: grp.Name}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, nil
	}

	return groups, nil
}

func (r *GroupRepository) PartiallyUpdateGroup(grp *group.Group) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := group.TableName
	val := fmt.Sprintf("%s=CASE WHEN $1 <> '' THEN $1 ELSE %s END", group.Name, group.Name)
	cnd := fmt.Sprintf("%s=$2", group.Id)
	query := fmt.Sprintf(template, tbl, val, cnd)

	rows, err := db.Query(query, grp.Name, grp.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *GroupRepository) DeleteGroup(id int) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.DELETE_FROM_TBL_WHERE_CND
	tbl := group.TableName
	cnd := fmt.Sprintf("%s=$1", group.Id)
	query := fmt.Sprintf(template, tbl, cnd)

	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *GroupRepository) IsGroupExists(idOrName interface{}) (bool, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	var template, col, tbl, cnd, query string
	var rows *sql.Rows
	var err error

	if reflect.TypeOf(idOrName) == reflect.TypeOf(0) {
		template = qconsts.SELECT_COL_FROM_TBL_WHERE_CND
		col = group.Id
		tbl = group.TableName
		cnd = fmt.Sprintf("%s=$1", group.Id)
		query = fmt.Sprintf(template, col, tbl, cnd)
		rows, err = db.Query(query, idOrName.(int))
	} else {
		template = qconsts.SELECT_COL_FROM_TBL_WHERE_CND
		col = group.Name
		tbl = group.TableName
		cnd = fmt.Sprintf("%s=$1", group.Name)
		query = fmt.Sprintf(template, col, tbl, cnd)
		rows, err = db.Query(query, idOrName.(string))
	}
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if isRowPresent := rows.Next(); !isRowPresent {
		return false, nil
	}

	return true, nil
}
