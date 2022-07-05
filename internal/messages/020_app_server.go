package messages

import (
	"fmt"
	"os"

	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func InfoServerStartedSuccessfullyAtLocalAddress(host string, port int) *logger.Log {
	return &logger.Log{Message: "Server was successfully started at local address: " + fmt.Sprintf("%s:%d", host, port)}
}

func InfoServerShutedDownCorrectly() *logger.Log {
	return &logger.Log{Message: "Server was correctly shuted down"}
}

func FatalFailureOnServerRunning(err error) *logger.Log {
	return &logger.Log{ErrCode: 20, Message: "Failure on server running. Error: " + err.Error(), ErrLevel: logger.ErrLevelFatal}
}

func WarningCannotGetLocalIP(err error) *logger.Log {
	return &logger.Log{ErrCode: 21, Message: "Cannot get local IP. Error: " + err.Error(), ErrLevel: logger.ErrLevelWarning}
}

func InfoRecivedSignal(signal os.Signal) *logger.Log {
	return &logger.Log{Message: "Recived signal: " + signal.String()}
}
