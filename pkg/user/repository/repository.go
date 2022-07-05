package repository

import (
	"database/sql"
	"fmt"
	"reflect"

	perm "github.com/mirumyantsev/video_hosting/internal/permission"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type UserRepository struct {
	cfg *config.Config
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	return &UserRepository{cfg: cfg}
}

func (r *UserRepository) CreateUser(usr *user.User) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s, %s, %s, %s, %s, %s, %s, %s)", user.TableName,
		user.Username, user.PasswordHash, user.IsActive, user.IsSuperuser, user.IsStaff,
		user.FirstName, user.LastName, user.JoiningDate, user.LastLogin)
	val := "($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	query := fmt.Sprintf(template, tbl, val)

	if _, err := db.Query(query, usr.Username, usr.PasswordHash, usr.IsActive, usr.IsSuperuser,
		usr.IsStaff, usr.FirstName, usr.LastName, usr.JoiningDate, usr.LastLogin); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(id int) (*user.User, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s, %s", user.Id, user.Username,
		user.PasswordHash, user.IsActive, user.IsSuperuser, user.IsStaff, user.FirstName,
		user.LastName, user.JoiningDate, user.LastLogin)
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1", user.Id)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var usr user.User
	if err := db.Get(&usr, query, id); err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *UserRepository) GetAllUsers(urlparams *user.Pagin) (map[int]*user.User, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.PAGINATION_COL_TBL_CND_PAG_TBL_PAG_LIM
	col := "*"
	tbl := user.TableName
	cnd := user.Id
	lim := urlparams.Limit
	pag := urlparams.Page
	query := fmt.Sprintf(template, col, tbl, cnd, pag, tbl, pag, lim)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = map[int]*user.User{}
	var usr user.User
	for rows.Next() {
		if err := rows.Scan(&usr.Id, &usr.Username, &usr.PasswordHash, &usr.IsActive, &usr.IsSuperuser,
			&usr.IsStaff, &usr.FirstName, &usr.LastName, &usr.JoiningDate, &usr.LastLogin); err != nil {
			return nil, err
		}
		users[usr.Id] = &user.User{Id: usr.Id, Username: usr.Username, PasswordHash: usr.PasswordHash,
			IsActive: usr.IsActive, IsSuperuser: usr.IsSuperuser, IsStaff: usr.IsStaff,
			FirstName: usr.FirstName, LastName: usr.LastName, JoiningDate: usr.JoiningDate,
			LastLogin: usr.LastLogin}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users, nil
}

func (r *UserRepository) UpdateUserPassword(namepass *auth.Namepass) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := user.TableName
	val := fmt.Sprintf("%s=CASE WHEN $1 <> '' THEN $1 ELSE %s END", auth.PasswordHash, auth.PasswordHash)
	cnd := fmt.Sprintf("%s=$2", auth.Username)
	query := fmt.Sprintf(template, tbl, val, cnd)

	rows, err := db.Query(query, namepass.PasswordHash, namepass.Username)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *UserRepository) PartiallyUpdateUser(usr *user.User) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := user.TableName
	val := fmt.Sprintf("%s=CASE WHEN $1 <> '' THEN $1 ELSE %s END, ", user.Username, user.Username) +
		fmt.Sprintf("%s=$2, ", user.IsActive) +
		fmt.Sprintf("%s=$3, ", user.IsSuperuser) +
		fmt.Sprintf("%s=$4, ", user.IsStaff) +
		fmt.Sprintf("%s=CASE WHEN $5 <> '' THEN $5 ELSE %s END, ", user.FirstName, user.FirstName) +
		fmt.Sprintf("%s=CASE WHEN $6 <> '' THEN $6 ELSE %s END", user.LastName, user.LastName)
	cnd := fmt.Sprintf("%s=$7", user.Id)
	query := fmt.Sprintf(template, tbl, val, cnd)

	rows, err := db.Query(query, usr.Username, usr.IsActive, usr.IsSuperuser, usr.IsStaff,
		usr.FirstName, usr.LastName, usr.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.DELETE_FROM_TBL_WHERE_CND
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1", user.Id)
	query := fmt.Sprintf(template, tbl, cnd)

	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *UserRepository) IsUserSuperuserOrStaff(username string) (bool, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s OR %s", user.IsSuperuser, user.IsStaff)
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1", user.Username)
	query := fmt.Sprintf(template, col, tbl, cnd)

	rows, err := db.Query(query, username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	queryRes := false
	for rows.Next() {
		if err := rows.Scan(&queryRes); err != nil {
			return false, err
		}
	}

	return queryRes, nil
}

func (r *UserRepository) IsUserHavePersonalPermission(userId int, userPerm string) (bool, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL1_FROM_TBL1_WHERE_CND1_SELECT_COL2_FROM_TBL2_CND2
	col1 := perm.Id
	tbl1 := perm.UPTableName
	cnd1 := fmt.Sprintf("%s=$1 AND %s", perm.UserId, perm.PermId)
	col2 := perm.Id
	tbl2 := "public.perms"
	cnd2 := "code_name=$2"
	query := fmt.Sprintf(template, col1, tbl1, cnd1, col2, tbl2, cnd2)

	rows, err := db.Query(query, userId, userPerm)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if isRowPresent := rows.Next(); !isRowPresent {
		return false, nil
	}

	return true, nil
}

func (r *UserRepository) IsUserExists(idOrUsername interface{}) (bool, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	var template, col, tbl, cnd, query string
	var rows *sql.Rows
	var err error

	if reflect.TypeOf(idOrUsername) == reflect.TypeOf(0) {
		template = qconsts.SELECT_COL_FROM_TBL_WHERE_CND
		col = user.Id
		tbl = user.TableName
		cnd = fmt.Sprintf("%s=$1", user.Id)
		query = fmt.Sprintf(template, col, tbl, cnd)
		rows, err = db.Query(query, idOrUsername.(int))
	} else {
		template = qconsts.SELECT_COL_FROM_TBL_WHERE_CND
		col = user.Username
		tbl = user.TableName
		cnd = fmt.Sprintf("%s=$1", user.Username)
		query = fmt.Sprintf(template, col, tbl, cnd)
		rows, err = db.Query(query, idOrUsername.(string))
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

func (r *UserRepository) GetUserId(username string) (int, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := user.Id
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1", user.Username)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var idPtr *int
	if err := db.Get(&idPtr, query, username); err != nil {
		return -1, err
	}

	return *idPtr, nil
}
