package usecase

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/info"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type InfoUseCase struct {
	infoRepo info.InfoRepository
}

func NewInfoUseCase(infoRepo info.InfoRepository) *InfoUseCase {
	return &InfoUseCase{
		infoRepo: infoRepo,
	}
}

func (u *InfoUseCase) CreateInfo(nfo *info.Info) error {
	return u.infoRepo.CreateInfo(nfo)
}

func (u *InfoUseCase) GetInfo(id int) (*info.Info, error) {
	return u.infoRepo.GetInfo(id)
}

func (u *InfoUseCase) GetAllInfos(urlparams *user.Pagin) (map[int]*info.Info, error) {
	return u.infoRepo.GetAllInfos(urlparams)
}

func (u *InfoUseCase) PartiallyUpdateInfo(nfo *info.Info) error {
	return u.infoRepo.PartiallyUpdateInfo(nfo)
}

func (u *InfoUseCase) DeleteInfo(id int) error {
	return u.infoRepo.DeleteInfo(id)
}

func (u *InfoUseCase) BindJSONInfo(ctx *gin.Context) (*info.Info, error) {
	var nfo info.Info
	if err := ctx.BindJSON(&nfo); err != nil {
		return &nfo, err
	}
	return &nfo, nil
}

func (u *InfoUseCase) IsRequiredEmpty(stream string) bool {
	if stream == "" {
		return true
	}
	return false
}

func (u *InfoUseCase) IsInfoExists(id int) (bool, error) {
	exists, err := u.infoRepo.IsInfoExists(id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *InfoUseCase) AtoiRequestedId(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return -1, err
	}
	return id, nil
}
