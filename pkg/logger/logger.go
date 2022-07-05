package logger

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/pkg/timedate"
)

const (
	tab    = "\t"
	indent = "    "
)

func Init(ctx *gin.Context) *Log {
	var log Log
	if ctx != nil {
		log.ClientIP = ctx.ClientIP()
		log.SessionOwner = "unauthorized"
		log.RequestMethod = ctx.Request.Method
		log.RequestPath = ctx.Request.URL.Path
	}
	log.CreationDate = timedate.GetTimestamp()
	return &log
}

func Complete(log1, log2 *Log) {
	if log2.ErrLevel != "" {
		log1.ErrLevel = log2.ErrLevel
	}
	if log2.ClientIP != "" {
		log1.ClientIP = log2.ClientIP
	}
	if log2.SessionOwner != "" {
		log1.SessionOwner = log2.SessionOwner
	}
	if log2.RequestMethod != "" {
		log1.RequestMethod = log2.RequestMethod
	}
	if log2.RequestPath != "" {
		log1.RequestPath = log2.RequestPath
	}
	if log2.StatusCode != 0 {
		log1.StatusCode = log2.StatusCode
	}
	if log2.ErrCode != 0 {
		log1.ErrCode = log2.ErrCode
	}
	if log2.Message != "" {
		log1.Message = log2.Message
	}
	if log2.CreationDate != "" {
		log1.CreationDate = log2.CreationDate
	}
}

func Finish(log *Log) {
	if log.ErrLevel == "" {
		log.ErrLevel = ErrLevelInfo
	}
	if log.CreationDate == "" {
		log.CreationDate = timedate.GetTimestamp()
	}
}

func Printc(ctx *gin.Context, messageLog *Log) {
	log := Init(ctx)
	Complete(log, messageLog)
	Print(log)
}

func Print(log *Log) {
	fmt.Println(parseLogLine(log))
}

func parseLogLine(log *Log) string {
	return parseErrLevel(log) + parseHttpLine(log) +
		parseErrorPrefix(log) + ParseMessage(log) +
		parseCreationDate(log)
}

func parseErrLevel(log *Log) string {
	if log.ErrLevel != "" {
		return log.ErrLevel + tab
	}
	return ErrLevelInfo + tab
}

func parseHttpLine(log *Log) string {
	if log.ClientIP != "" {
		return log.ClientIP + indent +
			log.SessionOwner + indent +
			log.RequestMethod + indent +
			log.RequestPath + indent +
			strconv.Itoa(log.StatusCode) + tab
	}
	return ""
}

func parseErrorPrefix(log *Log) string {
	if log.ErrCode != 0 {
		return "ErrCode: " + strconv.Itoa(log.ErrCode) + ". "
	}
	return ""
}

func ParseMessage(log *Log) string {
	msgType := fmt.Sprintf("%T", log.Message)
	if msgType == "string" {
		return log.Message.(string) + tab
	} else if msgType == "*user.User" {
		return "Got user" + tab
	} else if msgType == "map[int]*user.User" {
		return "Got all users" + tab
	} else if msgType == "*group.Group" {
		return "Got group" + tab
	} else if msgType == "map[int]*group.Group" {
		return "Got all groups" + tab
	} else if msgType == "*permission.PermIds" {
		return "Got user permissions" + tab
	} else if msgType == "map[int]*permission.Perm" {
		return "Got all permissions" + tab
	} else if msgType == "*info.Info" {
		return "Got info" + tab
	} else if msgType == "map[int]*info.Info" {
		return "Got all infos" + tab
	} else if msgType == "*video.Video" {
		return "Got video" + tab
	} else if msgType == "map[int]*video.Video" {
		return "Got all videos" + tab
	} else if msgType == "*group.GroupIds" {
		return "Got user groups" + tab
	} else if msgType == "*stream.Stream" {
		return "Got stream" + tab
	} else if msgType == "map[int]*stream.Stream" {
		return "Got all streams" + tab
	} else if msgType == "*download.Download" {
		return "Got download link" + tab
	}
	return "Got data of unknown type. Type: " + msgType + tab
}

func parseCreationDate(log *Log) string {
	if log.CreationDate != "" {
		return log.CreationDate
	}
	return timedate.GetTimestamp()
}
