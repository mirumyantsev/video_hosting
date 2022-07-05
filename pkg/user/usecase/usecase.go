package usecase

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	"github.com/mirumyantsev/video_hosting/pkg/hasher"
	"github.com/mirumyantsev/video_hosting/pkg/user"
)

type UserUseCase struct {
	cfg      *config.Config
	userRepo user.UserRepository
}

func NewUserUseCase(cfg *config.Config, userRepo user.UserRepository) *UserUseCase {
	return &UserUseCase{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (u *UserUseCase) CreateUser(usr *user.User) error {
	usr.IsActive = true
	usr.IsSuperuser = false
	usr.IsStaff = false
	usr.LastLogin = usr.JoiningDate
	return u.userRepo.CreateUser(usr)
}

func (u *UserUseCase) GetUser(id int) (*user.User, error) {
	return u.userRepo.GetUser(id)
}

func (u *UserUseCase) ParseURLParams(ctx *gin.Context) *user.Pagin {
	urlparams := ctx.Request.URL.Query()
	var pagin user.Pagin
	if lim := urlparams.Get("_limit"); lim != "" {
		pagin.Limit, _ = strconv.Atoi(lim)
	}
	if pg := urlparams.Get("_page"); pg != "" {
		pagin.Page, _ = strconv.Atoi(pg)
	}
	pagin.Page = pagin.Page*pagin.Limit - pagin.Limit
	if pagin.Limit == 0 {
		pagin.Limit = u.cfg.PaginationGetLimitDefault
	}
	return &pagin
}

func (u *UserUseCase) GetAllUsers(urlparams *user.Pagin) (map[int]*user.User, error) {
	return u.userRepo.GetAllUsers(urlparams)
}

func (u *UserUseCase) UpdateUserPassword(namepass *auth.Namepass) error {
	return u.userRepo.UpdateUserPassword(namepass)
}

func (u *UserUseCase) PartiallyUpdateUser(usr *user.User) error {
	return u.userRepo.PartiallyUpdateUser(usr)
}

func (u *UserUseCase) DeleteUser(id int) error {
	return u.userRepo.DeleteUser(id)
}

func (u *UserUseCase) IsUserSuperuserOrStaff(username string) (bool, error) {
	return u.userRepo.IsUserSuperuserOrStaff(username)
}

func (u *UserUseCase) IsUserHavePersonalPermission(userId int, userPerm string) (bool, error) {
	return u.userRepo.IsUserHavePersonalPermission(userId, userPerm)
}

func (u *UserUseCase) IsUserExists(idOrUsername interface{}) (bool, error) {
	exists, err := u.userRepo.IsUserExists(idOrUsername)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *UserUseCase) GetUserId(username string) (int, error) {
	return u.userRepo.GetUserId(username)
}

func (u *UserUseCase) BindJSONUser(ctx *gin.Context) (*user.User, error) {
	var usr user.User
	if err := ctx.BindJSON(&usr); err != nil {
		return &usr, err
	}
	if usr.PasswordHash != "" {
		usr.PasswordHash = hasher.GeneratePasswordHash(usr.PasswordHash, u.cfg.HashingPasswordSalt)
	}
	return &usr, nil
}

func (u *UserUseCase) IsRequiredEmpty(username, password string) bool {
	if username == "" || password == "" {
		return true
	}
	return false
}

func (u *UserUseCase) AtoiRequestedId(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return -1, err
	}
	return id, nil
}
