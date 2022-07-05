package repository

import (
	"fmt"

	perm "github.com/mirumyantsev/video_hosting/internal/permission"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type PermRepository struct {
	cfg *config.Config
}

func NewPermRepository(cfg *config.Config) *PermRepository {
	return &PermRepository{cfg: cfg}
}

func (r *PermRepository) GetAllPermissions(urlparams *user.Pagin) (map[int]*perm.Perm, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := "*"
	tbl := perm.TableName
	cnd := perm.Id
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms = map[int]*perm.Perm{}
	var prm perm.Perm
	for rows.Next() {
		if err := rows.Scan(&prm.Id, &prm.Name, &prm.Codename); err != nil {
			return nil, err
		}
		perms[prm.Id] = &perm.Perm{Id: prm.Id, Name: prm.Name, Codename: prm.Codename}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(perms) == 0 {
		return nil, nil
	}

	return perms, nil
}
