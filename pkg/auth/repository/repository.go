package repository

import (
	"fmt"

	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type AuthRepository struct {
	cfg *config.Config
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	return &AuthRepository{cfg: cfg}
}

func (r *AuthRepository) GetNamepass(namepass *auth.Namepass) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s", user.Username, user.PasswordHash)
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1 AND %s=$2", user.Username, user.PasswordHash)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var newNamepass auth.Namepass
	if err := db.Get(&newNamepass, query, namepass.Username, namepass.PasswordHash); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) UpdateUserPassword(namepass *auth.Namepass) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := user.TableName
	val := fmt.Sprintf("%s=CASE WHEN $1 <> '' THEN $1 ELSE %s END", user.PasswordHash, user.PasswordHash)
	cnd := fmt.Sprintf("%s=$2", user.Username)
	query := fmt.Sprintf(template, tbl, val, cnd)

	rows, err := db.Query(query, namepass.PasswordHash, namepass.Username)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *AuthRepository) IsUsernameAndPasswordExists(usename, passwordHash string) (bool, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := user.Id
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1 AND %s=$2", user.Username, user.PasswordHash)
	query := fmt.Sprintf(template, col, tbl, cnd)

	rows, err := db.Query(query, usename, passwordHash)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	isRowPresent := rows.Next()
	if !isRowPresent {
		return false, nil
	}

	return true, nil
}

func (r *AuthRepository) UpdateNamepassLastLogin(username, timestamp string) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := user.TableName
	val := fmt.Sprintf("%s=$1", user.LastLogin)
	cnd := fmt.Sprintf("%s=$2", user.Username)
	query := fmt.Sprintf(template, tbl, val, cnd)

	rows, err := db.Query(query, timestamp, username)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
