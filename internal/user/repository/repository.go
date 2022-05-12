package repository

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/mikerumy/vhosting/internal/user"
	"github.com/mikerumy/vhosting/pkg/config_tool"
	query_consts "github.com/mikerumy/vhosting/pkg/constants/query"
	"github.com/mikerumy/vhosting/pkg/db_tool"
)

type UserRepository struct {
	cfg config_tool.Config
}

func NewUserRepository(cfg config_tool.Config) *UserRepository {
	return &UserRepository{cfg: cfg}
}

func (r *UserRepository) CreateUser(usr user.User) error {
	db := db_tool.NewDBConnection(r.cfg)
	defer db_tool.CloseDBConnection(r.cfg, db)

	template := query_consts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s, %s, %s, %s, %s, %s, %s, %s)", user.TableName, user.Username,
		user.PasswordHash, user.IsActive, user.IsSuperUser, user.IsStaff, user.FirstName,
		user.LastName, user.JoiningDate, user.LastLogin)
	val := "($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	query := fmt.Sprintf(template, tbl, val)

	_, err := db.Query(query, usr.Username, usr.PasswordHash, usr.IsActive, usr.IsSuperUser,
		usr.IsStaff, usr.FirstName, usr.LastName, usr.JoiningDate, usr.LastLogin)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(id int) (*user.User, error) {
	db := db_tool.NewDBConnection(r.cfg)
	defer db_tool.CloseDBConnection(r.cfg, db)

	template := query_consts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s, %s", user.Id, user.Username,
		user.PasswordHash, user.IsActive, user.IsSuperUser, user.IsStaff, user.FirstName,
		user.LastName, user.JoiningDate, user.LastLogin)
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1", user.Id)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var usr user.User
	err := db.Get(&usr, query, id)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *UserRepository) GetAllUsers() (map[int]*user.User, error) {
	db := db_tool.NewDBConnection(r.cfg)
	defer db_tool.CloseDBConnection(r.cfg, db)

	template := query_consts.SELECT_COL_FROM_TBL
	col := "*"
	tbl := user.TableName
	query := fmt.Sprintf(template, col, tbl)

	var rows *sql.Rows
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = map[int]*user.User{}
	var usr user.User
	for rows.Next() {
		err = rows.Scan(&usr.Id, &usr.Username, &usr.PasswordHash, &usr.IsActive, &usr.IsSuperUser,
			&usr.IsStaff, &usr.FirstName, &usr.LastName, &usr.JoiningDate, &usr.LastLogin)
		if err != nil {
			return nil, err
		}
		users[usr.Id] = &user.User{Id: usr.Id, Username: usr.Username, PasswordHash: usr.PasswordHash,
			IsActive: usr.IsActive, IsSuperUser: usr.IsSuperUser, IsStaff: usr.IsStaff,
			FirstName: usr.FirstName, LastName: usr.LastName, JoiningDate: usr.JoiningDate,
			LastLogin: usr.LastLogin}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users, nil
}

func (r *UserRepository) PartiallyUpdateUser(usr *user.User) error {
	db := db_tool.NewDBConnection(r.cfg)
	defer db_tool.CloseDBConnection(r.cfg, db)

	template := query_consts.UPDATE_TBL_SET_VAL_WHERE_CND
	tbl := user.TableName
	val := fmt.Sprintf("%s=CASE WHEN $1 <> '' THEN $1 ELSE %s END, ", user.Username, user.Username) +
		fmt.Sprintf("%s=CASE WHEN $2 <> '' THEN $2 ELSE %s END, ", user.PasswordHash, user.PasswordHash) +
		fmt.Sprintf("%s=$3, ", user.IsActive) +
		fmt.Sprintf("%s=$4, ", user.IsSuperUser) +
		fmt.Sprintf("%s=$5, ", user.IsStaff) +
		fmt.Sprintf("%s=CASE WHEN $6 <> '' THEN $6 ELSE %s END, ", user.FirstName, user.FirstName) +
		fmt.Sprintf("%s=CASE WHEN $7 <> '' THEN $7 ELSE %s END", user.LastName, user.LastName)
	cnd := fmt.Sprintf("%s=$8", user.Id)
	query := fmt.Sprintf(template, tbl, val, cnd)

	var rows *sql.Rows
	rows, err := db.Query(query, usr.Username, usr.PasswordHash, usr.IsActive, usr.IsSuperUser,
		usr.IsStaff, usr.FirstName, usr.LastName, usr.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	db := db_tool.NewDBConnection(r.cfg)
	defer db_tool.CloseDBConnection(r.cfg, db)

	template := query_consts.DELETE_FROM_TBL_WHERE_CND
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1", user.Id)
	query := fmt.Sprintf(template, tbl, cnd)

	var rows *sql.Rows
	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *UserRepository) IsUserExists(idOrUsername interface{}) (bool, error) {
	db := db_tool.NewDBConnection(r.cfg)
	defer db_tool.CloseDBConnection(r.cfg, db)

	var template, col, tbl, cnd, query string
	var rows *sql.Rows
	var err error

	if reflect.TypeOf(idOrUsername) == reflect.TypeOf(0) {
		template = query_consts.SELECT_COL_FROM_TBL_WHERE_CND
		col = user.Id
		tbl = user.TableName
		cnd = fmt.Sprintf("%s=$1", user.Id)
		query = fmt.Sprintf(template, col, tbl, cnd)
		rows, err = db.Query(query, idOrUsername.(int))
	} else {
		template = query_consts.SELECT_COL_FROM_TBL_WHERE_CND
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

	rowIsPresent := rows.Next()
	if !rowIsPresent {
		return false, nil
	}

	return true, nil
}

func (r *UserRepository) GetUserId(username string) (int, error) {
	db := db_tool.NewDBConnection(r.cfg)
	defer db_tool.CloseDBConnection(r.cfg, db)

	template := query_consts.SELECT_COL_FROM_TBL_WHERE_CND
	col := user.Id
	tbl := user.TableName
	cnd := fmt.Sprintf("%s=$1", user.Username)
	query := fmt.Sprintf(template, col, tbl, cnd)

	var idPtr *int
	err := db.Get(&idPtr, query, username)
	if err != nil {
		return -1, err
	}

	return *idPtr, nil
}