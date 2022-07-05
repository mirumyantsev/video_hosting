package db_connect

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func CreateLocalDBConnection(cfg *config.Config) *sqlx.DB {
	var dbcfg DBConfig
	dbcfg.DBConnectionLatencyMilliseconds = cfg.DBConnectionLatencyMilliseconds
	dbcfg.DBConnectionShowStatus = cfg.DBConnectionShowStatus
	dbcfg.DBConnectionTimeoutSeconds = cfg.DBConnectionTimeoutSeconds
	dbcfg.DBHost = cfg.DBHost
	dbcfg.DBName = cfg.DBName
	dbcfg.DBPort = cfg.DBPort
	dbcfg.DBSSLEnable = cfg.DBSSLEnable
	dbcfg.DBUsername = cfg.DBUsername
	dbcfg.DBDriver = cfg.DBDriver
	dbcfg.DBPassword = cfg.DBPassword
	return connectToDB(&dbcfg)
}

func CreateOuterDBConnection(cfg *config.Config) *sqlx.DB {
	var dbcfg DBConfig
	dbcfg.DBConnectionLatencyMilliseconds = cfg.DBOConnectionLatencyMilliseconds
	dbcfg.DBConnectionShowStatus = cfg.DBOConnectionShowStatus
	dbcfg.DBConnectionTimeoutSeconds = cfg.DBOConnectionTimeoutSeconds
	dbcfg.DBHost = cfg.DBOHost
	dbcfg.DBName = cfg.DBOName
	dbcfg.DBPort = cfg.DBOPort
	dbcfg.DBSSLEnable = cfg.DBOSSLEnable
	dbcfg.DBUsername = cfg.DBOUsername
	dbcfg.DBDriver = cfg.DBODriver
	dbcfg.DBPassword = cfg.DBOPassword
	return connectToDB(&dbcfg)
}

func connectToDB(dbcfg *DBConfig) *sqlx.DB {
	timeAtStarting := time.Now()
	var db *sqlx.DB
	sslMode := "disable"
	if dbcfg.DBSSLEnable {
		sslMode = "enable"
	}
	go func() *sqlx.DB {
		for {
			db, _ = sqlx.Open(dbcfg.DBDriver, fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
				dbcfg.DBHost, dbcfg.DBPort, dbcfg.DBUsername, dbcfg.DBName, dbcfg.DBPassword, sslMode))
			time.Sleep(3 * time.Millisecond)
			if db.Ping() == nil {
				return db
			}
		}
	}()
	connLatency := time.Duration(dbcfg.DBConnectionLatencyMilliseconds)
	time.Sleep(connLatency * time.Millisecond)
	connTimeout := dbcfg.DBConnectionTimeoutSeconds
	for t := connTimeout; t > 0; t-- {
		if db != nil {
			if dbcfg.DBConnectionShowStatus {
				logger.Print(msg.InfoEstablishedOpenedDBConnection(timeAtStarting))
				return db
			}
			return db
		}
		time.Sleep(time.Second)
	}
	logger.Print(msg.ErrorTimeWaitingOfDBConnectionExceededLimit(connTimeout))
	return nil
}

func CloseDBConnection(cfg *config.Config, db *sqlx.DB) {
	if err := db.Close(); err != nil {
		logger.Print(msg.ErrorCannotCloseDBConnection(err))
		return
	}
	if cfg.DBConnectionShowStatus {
		logger.Print(msg.InfoEstablishedClosedConnectionToDB())
		return
	}
}
