package repository

import (
	"fmt"

	sess "github.com/mirumyantsev/video_hosting/internal/session"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
)

type SessRepository struct {
	cfg *config.Config
}

func NewSessRepository(cfg *config.Config) *SessRepository {
	return &SessRepository{cfg: cfg}
}

func (r *SessRepository) CreateSession(session *sess.Session) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s)", sess.TableName, sess.Content, sess.CreationDate)
	val := "($1, $2)"
	query := fmt.Sprintf(template, tbl, val)

	_, err := db.Query(query, session.Content, session.CreationDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessRepository) GetSessionAndDate(token string) (*sess.Session, error) {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.SELECT_COL_FROM_TBL_WHERE_CND
	col := fmt.Sprintf("%s, %s", sess.Content, sess.CreationDate)
	tbl := sess.TableName
	cnd := fmt.Sprintf("%s='%s'", sess.Content, token)
	query := fmt.Sprintf(template, col, tbl, cnd)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var session sess.Session
	for rows.Next() {
		if err := rows.Scan(&session.Content, &session.CreationDate); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessRepository) DeleteSession(token string) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.DELETE_FROM_TBL_WHERE_CND
	tbl := sess.TableName
	cnd := fmt.Sprintf("%s=$1", sess.Content)
	query := fmt.Sprintf(template, tbl, cnd)

	rows, err := db.Query(query, token)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
