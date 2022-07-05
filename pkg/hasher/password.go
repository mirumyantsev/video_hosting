package hasher

import (
	"crypto/sha256"
	"fmt"

	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
)

func GeneratePasswordHash(password, salt string) string {
	if password == "" {
		return ""
	}
	hash := sha256.New()
	if _, err := hash.Write([]byte(password)); err != nil {
		logger.Print(msg.ErrorCannotWriteBytesIntoInternalVariable(err))
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
