package logger

type Log struct {
	Id            int         `db:"id"`
	ErrLevel      string      `db:"error_level"` // "info", "warning", "error", "fatal"
	ClientIP      string      `db:"client_ip"`
	SessionOwner  string      `db:"session_owner"`
	RequestMethod string      `db:"request_method"` // "POST", "GET", "PATCH", "DELETE"
	RequestPath   string      `db:"request_path"`
	StatusCode    int         `db:"status_code"`
	ErrCode       int         `db:"error_code"`
	Message       interface{} `db:"message"`
	CreationDate  string      `db:"creation_date"`
}
