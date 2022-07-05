package messages

import (
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func InfoYouHaveSuccessfullySignedOut() *logger.Log {
	return &logger.Log{StatusCode: 202, Message: "You have successfully signed-out"}
}

func ErrorUserWithEnteredUsernameOrPasswordIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 100, Message: "User with entered username or password is not exist", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotGenerateToken(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 101, Message: "Cannot generate token. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoYouHaveSuccessfullySignedIn() *logger.Log {
	return &logger.Log{StatusCode: 202, Message: "You have successfully signed-in"}
}

func ErrorYouMustBeSignedInForChangingPassword() *logger.Log {
	return &logger.Log{StatusCode: 401, ErrCode: 102, Message: "You must be signed-in for changing password", ErrLevel: logger.ErrLevelError}
}

func ErrorCannotParseToken(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 103, Message: "Cannot parse token. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func ErrorUserWithSuchUsernameOrPasswordIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 104, Message: "User with such username or password is not exist", ErrLevel: logger.ErrLevelError}
}

func ErrorPasswordCannotBeEmpty() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 105, Message: "Password cannot be empty"}
}

func ErrorCannotUpdateUserPassword(err error) *logger.Log {
	return &logger.Log{StatusCode: 500, ErrCode: 106, Message: "Cannot update user password. Error: " + err.Error(), ErrLevel: logger.ErrLevelError}
}

func InfoYouHaveSuccessfullyChangedPassword() *logger.Log {
	return &logger.Log{StatusCode: 202, Message: "You have successfully changed password"}
}

func ErrorYouMustBeSignedInForSigningOut() *logger.Log {
	return &logger.Log{StatusCode: 401, ErrCode: 107, Message: "You must be signed-in for signing-out", ErrLevel: logger.ErrLevelError}
}

func ErrorUserWithThisUsernameIsNotExist() *logger.Log {
	return &logger.Log{StatusCode: 400, ErrCode: 108, Message: "User with this username is not exist", ErrLevel: logger.ErrLevelError}
}
