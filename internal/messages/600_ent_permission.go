package messages

import (
	perm "github.com/mirumyantsev/video_hosting/internal/permission"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func ErrorCannotGetAllPerms(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 600, Message: "Cannot get all permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoNoPermsAvailable() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "No permissions available"}
}

func InfoGotAllPerms(groups map[int]*perm.Perm) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: groups}
}

func ErrorPermIdsCannotBeEmpty() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 601, Message: "Permission IDs cannot be empty", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotSetUserPerms(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 602, Message: "Cannot set user permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoUserPermsSet() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "User permissions set"}
}

func ErrorCannotGetUserPerms(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 603, Message: "Cannot get user permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGotUserPerms(permIds *perm.PermIds) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: permIds}
}

func ErrorCannotDeleteUserPerms(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 604, Message: "Cannot delete user permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoUserPermsDeleted() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "User permissions deleted"}
}

func ErrorCannotSetGroupPerms(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 605, Message: "Cannot set group permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGroupPermsSet() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Group permissions set"}
}

func ErrorCannotGetGroupPerms(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 606, Message: "Cannot get group permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGotGroupPerms(perms *perm.PermIds) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: perms}
}

func ErrorCannotDeleteGroupPerms(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 607, Message: "Cannot delete group permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGroupPermsDeleted() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "Group permissions deleted"}
}
