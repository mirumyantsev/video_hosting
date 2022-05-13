package messages

import (
	"fmt"

	lg "github.com/mikerumy/vhosting/internal/logging"
	"github.com/mikerumy/vhosting/pkg/logger"
)

func InfoYouHaveSuccessfullySignedOut() *lg.Log {
	return &lg.Log{StatusCode: 202, Message: "You have successfully signed-out.", ErrorLevel: logger.ErrLevelInfo}
}

func ErrorUserWithEnteredUsernameOrPasswordIsNotExist() *lg.Log {
	return &lg.Log{StatusCode: 400, ErrorCode: 100, Message: "User with entered username or password is not exist.", ErrorLevel: logger.ErrLevelError}
}

func ErrorCannotGenerateToken(err error) *lg.Log {
	return &lg.Log{StatusCode: 500, ErrorCode: 101, Message: fmt.Sprintf("Cannot generate token. Error: %s.", err.Error()), ErrorLevel: logger.ErrLevelError}
}

func InfoYouHaveSuccessfullySignedIn() *lg.Log {
	return &lg.Log{StatusCode: 202, Message: "You have successfully signed-in.", ErrorLevel: logger.ErrLevelInfo}
}

func ErrorYouMustBeSignedInForChangingPassword() *lg.Log {
	return &lg.Log{StatusCode: 401, ErrorCode: 102, Message: "You must be signed-in for changing password.", ErrorLevel: logger.ErrLevelError}
}

func ErrorCannotParseToken(err error) *lg.Log {
	return &lg.Log{StatusCode: 500, ErrorCode: 103, Message: fmt.Sprintf("Cannot parse token. Error: %s.", err.Error()), ErrorLevel: logger.ErrLevelError}
}

func ErrorUserWithEnteredUsernameIsNotExist() *lg.Log {
	return &lg.Log{StatusCode: 400, ErrorCode: 104, Message: "User with entered username is not exist.", ErrorLevel: logger.ErrLevelError}
}

func ErrorEnteredUsernameIsIncorrect() *lg.Log {
	return &lg.Log{StatusCode: 400, ErrorCode: 105, Message: "Entered username is incorrect.", ErrorLevel: logger.ErrLevelError}
}

func ErrorCannotUpdateNamepassPassword(err error) *lg.Log {
	return &lg.Log{StatusCode: 500, ErrorCode: 106, Message: fmt.Sprintf("Cannot update user password. Error: %s.", err.Error()), ErrorLevel: logger.ErrLevelError}
}

func InfoYouHaveSuccessfullyChangedPassword() *lg.Log {
	return &lg.Log{StatusCode: 202, Message: "You have successfully changed password.", ErrorLevel: logger.ErrLevelInfo}
}

func ErrorYouMustBeSignedIn() *lg.Log {
	return &lg.Log{StatusCode: 400, ErrorCode: 107, Message: "You must be signed-in.", ErrorLevel: logger.ErrLevelError}
}

func ErrorYouMustBeSignedInForSigningOut() *lg.Log {
	return &lg.Log{StatusCode: 401, ErrorCode: 108, Message: "You must be signed-in for signing-out.", ErrorLevel: logger.ErrLevelError}
}

func ErrorUserWithThisUsernameIsNotExist() *lg.Log {
	return &lg.Log{StatusCode: 400, ErrorCode: 109, Message: "User with this username is not exist.", ErrorLevel: logger.ErrLevelError}
}