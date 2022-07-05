package messages

import (
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

func ErrorCannotBindInputData(err error) *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 200, Message: "Cannot bind input data. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorUsernameAndPasswordCannotBeEmpty() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 201, Message: "Username and password cannot be empty", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCheckUserExistence(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 202, Message: "Cannot check user existence. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorUserWithEnteredUsernameIsExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 203, Message: "User with entered username is exist", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCreateUser(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 204, Message: "Cannot create user. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoUserCreated() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "User created"}
}

func ErrorCannotConvertRequestedIDToTypeInt(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 205, Message: "Cannot convert requested ID to type int. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorUserWithRequestedIDIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 206, Message: "User with requested ID is not exist", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotGetUser(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 207, Message: "Cannot get user. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoGotUser(usr *user.User) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: usr}
}

func ErrorCannotGetAllUsers(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 208, Message: "Cannot get all users. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoNoUsersAvailable() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "No users available"}
}

func InfoGotAllUsers(users map[int]*user.User) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: users}
}

func ErrorCannotPartiallyUpdateUser(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 209, Message: "Cannot partially update user. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoUserPartiallyUpdated() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "User partially updated"}
}

func ErrorCannotDeleteUser(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 210, Message: "Cannot delete user. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoUserDeleted() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "User deleted"}
}

func ErrorYouHaveNotEnoughPermissions() *logger.Log {
	return &logger.Log{StatusCode: 403, ErrCode: 211, Message: "You have not enough permissions", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCheckSuperuserStaffPermissions(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 212, Message: "Cannot check superuser/staff permissions. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorCannotCheckPersonalPermission(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 213, Message: "Cannot check personal permission. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorUserWithEnteredUsernameIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 214, Message: "User with entered username is not exist", ErrLevel: logger.ErrLevelError}
}

func InfoUserPasswordChanged() *logger.Log {
	return &logger.Log{StatusCode: 200, Message: "User password changed"}
}
