package usecase

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/video"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type VideoUseCase struct {
	videoRepo video.VideoRepository
}

func NewVideoUseCase(videoRepo video.VideoRepository) *VideoUseCase {
	return &VideoUseCase{
		videoRepo: videoRepo,
	}
}

func (u *VideoUseCase) CreateVideo(nfo *video.Video) error {
	return u.videoRepo.CreateVideo(nfo)
}

func (u *VideoUseCase) GetVideo(id int) (*video.Video, error) {
	return u.videoRepo.GetVideo(id)
}

func (u *VideoUseCase) GetAllVideos(urlparams *user.Pagin) (map[int]*video.Video, error) {
	return u.videoRepo.GetAllVideos(urlparams)
}

func (u *VideoUseCase) PartiallyUpdateVideo(nfo *video.Video) error {
	return u.videoRepo.PartiallyUpdateVideo(nfo)
}

func (u *VideoUseCase) DeleteVideo(id int) error {
	return u.videoRepo.DeleteVideo(id)
}

func (u *VideoUseCase) BindJSONVideo(ctx *gin.Context) (*video.Video, error) {
	var vid video.Video
	if err := ctx.BindJSON(&vid); err != nil {
		return &vid, err
	}
	return &vid, nil
}

func (u *VideoUseCase) IsRequiredEmpty(url, filename string) bool {
	if url == "" || filename == "" {
		return true
	}
	return false
}

func (u *VideoUseCase) IsVideoExists(id int) (bool, error) {
	exists, err := u.videoRepo.IsVideoExists(id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *VideoUseCase) AtoiRequestedId(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return -1, err
	}
	return id, nil
}
