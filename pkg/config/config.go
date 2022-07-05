package config

import (
	"os"
	"strconv"
	"strings"

	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*Config, error) {
	// Parse config file path
	path = path[:len(path)-4]
	lastDirIndex := strings.LastIndex(path, "/")
	viper.AddConfigPath(path[:lastDirIndex+1])
	viper.SetConfigName(path[lastDirIndex+1:])

	// Load data from config file
	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return &cfg, err
	}

	param := ""

	param = "DB_CONNECTION_LATENCY_MILLISECONDS"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 50
		cfg.DBConnectionLatencyMilliseconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBConnectionLatencyMilliseconds = val
	}

	param = "DB_CONNECTION_SHOW_STATUS"
	if val, err := strconv.ParseBool(os.Getenv(param)); err != nil {
		defaultVal := false
		cfg.DBConnectionShowStatus = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBConnectionShowStatus = val
	}

	param = "DB_CONNECTION_TIMEOUT_SECONDS"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 5
		cfg.DBConnectionTimeoutSeconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBConnectionTimeoutSeconds = val
	}

	param = "DB_DRIVER"
	if os.Getenv(param) == "" {
		defaultVal := "mysql"
		cfg.DBDriver = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBDriver = os.Getenv(param)
	}

	param = "DB_HOST"
	if os.Getenv(param) == "" {
		defaultVal := "127.0.0.1"
		cfg.DBHost = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBHost = os.Getenv(param)
	}

	param = "DB_NAME"
	if os.Getenv(param) == "" {
		defaultVal := "vhs"
		cfg.DBName = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBName = os.Getenv(param)
	}

	param = "DB_PASSWORD"
	if os.Getenv(param) == "" {
		defaultVal := "1234"
		cfg.DBPassword = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBPassword = os.Getenv(param)
	}

	param = "DB_PORT"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 2345
		cfg.DBPort = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBPort = val
	}

	param = "DB_SSL_ENABLE"
	if val, err := strconv.ParseBool(os.Getenv(param)); err != nil {
		defaultVal := false
		cfg.DBSSLEnable = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBSSLEnable = val
	}

	param = "DB_USERNAME"
	if os.Getenv(param) == "" {
		defaultVal := "joe"
		cfg.DBUsername = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBUsername = os.Getenv(param)
	}

	param = "DBO_CONNECTION_LATENCY_MILLISECONDS"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 500
		cfg.DBOConnectionLatencyMilliseconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOConnectionLatencyMilliseconds = val
	}

	param = "DBO_CONNECTION_SHOW_STATUS"
	if val, err := strconv.ParseBool(os.Getenv(param)); err != nil {
		defaultVal := false
		cfg.DBOConnectionShowStatus = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOConnectionShowStatus = val
	}

	param = "DBO_CONNECTION_TIMEOUT_SECONDS"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 5
		cfg.DBOConnectionTimeoutSeconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOConnectionTimeoutSeconds = val
	}

	param = "DBO_DRIVER"
	if os.Getenv(param) == "" {
		defaultVal := "mysql"
		cfg.DBODriver = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBODriver = os.Getenv(param)
	}

	param = "DBO_HOST"
	if os.Getenv(param) == "" {
		defaultVal := "127.0.0.1"
		cfg.DBOHost = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOHost = os.Getenv(param)
	}

	param = "DBO_NAME"
	if os.Getenv(param) == "" {
		defaultVal := "l333"
		cfg.DBOName = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOName = os.Getenv(param)
	}

	param = "DBO_PASSWORD"
	if os.Getenv(param) == "" {
		defaultVal := "1234"
		cfg.DBOPassword = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOPassword = os.Getenv(param)
	}

	param = "DBO_PORT"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 2345
		cfg.DBOPort = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOPort = val
	}

	param = "DBO_SSL_ENABLE"
	if val, err := strconv.ParseBool(os.Getenv(param)); err != nil {
		defaultVal := false
		cfg.DBOSSLEnable = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOSSLEnable = val
	}

	param = "DBO_USERNAME"
	if os.Getenv(param) == "" {
		defaultVal := "john"
		cfg.DBOUsername = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.DBOUsername = os.Getenv(param)
	}

	param = "HASHING_PASSWORD_SALT"
	if os.Getenv(param) == "" {
		defaultVal := "SdD2Sdf@dFhSe#r"
		cfg.HashingPasswordSalt = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.HashingPasswordSalt = os.Getenv(param)
	}

	param = "HASHING_TOKEN_SIGNING_KEY"
	if os.Getenv(param) == "" {
		defaultVal := "gHs@dHk4Bs#v5HeK4h"
		cfg.HashingTokenSigningKey = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.HashingTokenSigningKey = os.Getenv(param)
	}

	param = "SERVER_HOST"
	if os.Getenv(param) == "" {
		defaultVal := "localhost"
		cfg.ServerHost = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.ServerHost = os.Getenv(param)
	}

	param = "SERVER_PORT"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 8000
		cfg.ServerPort = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.ServerPort = val
	}

	param = "SERVER_READ_TIMEOUT_SECONDS"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 15
		cfg.ServerReadTimeoutSeconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.ServerReadTimeoutSeconds = val
	}

	param = "SERVER_WRITE_TIMEOUT_SECONDS"
	if val, err := strconv.Atoi(os.Getenv(param)); err != nil {
		defaultVal := 15
		cfg.ServerWriteTimeoutSeconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.ServerWriteTimeoutSeconds = val
	}

	param = "pagination.getLimitDefault"
	if val := viper.GetInt(param); val == 0 {
		defaultVal := 30
		cfg.PaginationGetLimitDefault = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.PaginationGetLimitDefault = val
	}

	cfg.ServerDebugEnable = viper.GetBool("server.debugEnable")

	param = "server.maxHeaderBytes"
	if val := viper.GetInt(param); val == 0 {
		defaultVal := 1048576 // 1 megabyte
		cfg.ServerMaxHeaderBytes = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.ServerMaxHeaderBytes = val
	}

	param = "session.ttlHours"
	if val := viper.GetInt(param); val == 0 {
		defaultVal := 168 // 7 days
		cfg.SessionTTLHours = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.SessionTTLHours = val
	}

	param = "stream.iceServers"
	if val := viper.GetStringSlice(param); len(val) == 0 {
		defaultICEServer := "stun:stun.l.google.com:19302"
		defaultVal := []string{defaultICEServer}
		cfg.StreamICEServers = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, "[\""+defaultICEServer+"\"]"))
	} else {
		cfg.StreamICEServers = val
	}

	param = "stream.snapshotPeriodSeconds"
	if val := viper.GetInt(param); val == 0 {
		defaultVal := 60
		cfg.StreamSnapshotPeriodSeconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.StreamSnapshotPeriodSeconds = val
	}

	cfg.StreamSnapshotShowStatus = viper.GetBool("stream.snapshotShowStatus")
	cfg.StreamSnapshotsEnable = viper.GetBool("stream.snapshotsEnable")

	param = "stream.streamsUpdatePeriodSeconds"
	if val := viper.GetInt(param); val == 0 {
		defaultVal := 60
		cfg.StreamStreamsUpdatePeriodSeconds = defaultVal
		logger.Print(msg.WarningCannotConvertCvar(param, defaultVal))
	} else {
		cfg.StreamStreamsUpdatePeriodSeconds = val
	}

	return &cfg, nil
}
