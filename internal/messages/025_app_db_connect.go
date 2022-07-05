package messages

import (
	"strconv"
	"time"

	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func InfoEstablishedOpenedDBConnection(timeSinceOpen time.Time) *logger.Log {
	return &logger.Log{Message: "Established opening of connection to DB in %s" + time.Since(timeSinceOpen).Round(time.Millisecond).String()}
}

func ErrorTimeWaitingOfDBConnectionExceededLimit(timeout int) *logger.Log {
	return &logger.Log{ErrCode: 25, Message: "Time waiting of DB connection exceeded limit (" + strconv.Itoa(timeout) + " seconds)", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCloseDBConnection(err error) *logger.Log {
	return &logger.Log{ErrCode: 26, Message: "Cannot close DB connection. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoEstablishedClosedConnectionToDB() *logger.Log {
	return &logger.Log{Message: "Established closing of connection to DB"}
}
