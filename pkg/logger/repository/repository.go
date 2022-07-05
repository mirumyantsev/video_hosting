package repository

import (
	"fmt"

	"github.com/mirumyantsev/video_hosting/pkg/config"
	qconsts "github.com/mirumyantsev/video_hosting/pkg/constants/query"
	"github.com/mirumyantsev/video_hosting/pkg/db_connect"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

type LogRepository struct {
	cfg *config.Config
}

func NewLogRepository(cfg *config.Config) *LogRepository {
	return &LogRepository{cfg: cfg}
}

func (r *LogRepository) CreateLogRecord(log *logger.Log) error {
	db := db_connect.CreateLocalDBConnection(r.cfg)
	defer db_connect.CloseDBConnection(r.cfg, db)

	template := qconsts.INSERT_INTO_TBL_VALUES_VAL
	tbl := fmt.Sprintf("%s (%s, %s, %s, %s, %s, %s, %s, %s, %s)", logger.TableName,
		logger.ErrLevel, logger.ClientID, logger.SessionOwner, logger.RequestMethod, logger.RequestPath,
		logger.StatusCode, logger.ErrCode, logger.Message, logger.CreationDate)
	val := "($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	query := fmt.Sprintf(template, tbl, val)

	if _, err := db.Query(query, log.ErrLevel, log.ClientIP, log.SessionOwner,
		log.RequestMethod, log.RequestPath, log.StatusCode, log.ErrCode,
		log.Message.(string), log.CreationDate); err != nil {
		return err
	}

	return nil
}
