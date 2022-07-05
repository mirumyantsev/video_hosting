package usecase

import (
	"fmt"
	"strings"

	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/download"
)

type DownloadUseCase struct {
	cfg *config.Config
}

func NewDownloadUseCase(cfg *config.Config) *DownloadUseCase {
	return &DownloadUseCase{
		cfg: cfg,
	}
}

func (u *DownloadUseCase) IsValidExtension(file_name string) bool {
	extension := file_name[len(file_name)-4:]
	strings.ToLower(extension)
	if extension != ".mp4" {
		return false
	}
	return true
}

func (u *DownloadUseCase) CreateDownloadLink(local_file_path string) *download.Download {
	var dload download.Download
	dload.DownloadLink = fmt.Sprintf("http://%s:%d/media/%s", u.cfg.ServerIP,
		u.cfg.ServerPort, local_file_path)
	return &dload
}
