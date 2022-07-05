package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirumyantsev/video_hosting/internal/group"
	grouphandler "github.com/mirumyantsev/video_hosting/internal/group/handler"
	grouprepo "github.com/mirumyantsev/video_hosting/internal/group/repository"
	groupusecase "github.com/mirumyantsev/video_hosting/internal/group/usecase"
	"github.com/mirumyantsev/video_hosting/internal/info"
	infohandler "github.com/mirumyantsev/video_hosting/internal/info/handler"
	inforepo "github.com/mirumyantsev/video_hosting/internal/info/repository"
	infousecase "github.com/mirumyantsev/video_hosting/internal/info/usecase"
	msg "github.com/mirumyantsev/video_hosting/internal/messages"
	perm "github.com/mirumyantsev/video_hosting/internal/permission"
	permhandler "github.com/mirumyantsev/video_hosting/internal/permission/handler"
	permrepo "github.com/mirumyantsev/video_hosting/internal/permission/repository"
	permusecase "github.com/mirumyantsev/video_hosting/internal/permission/usecase"
	sess "github.com/mirumyantsev/video_hosting/internal/session"
	sessrepo "github.com/mirumyantsev/video_hosting/internal/session/repository"
	sessusecase "github.com/mirumyantsev/video_hosting/internal/session/usecase"
	"github.com/mirumyantsev/video_hosting/internal/video"
	videohandler "github.com/mirumyantsev/video_hosting/internal/video/handler"
	videorepo "github.com/mirumyantsev/video_hosting/internal/video/repository"
	videousecase "github.com/mirumyantsev/video_hosting/internal/video/usecase"
	"github.com/mirumyantsev/video_hosting/pkg/auth"
	authhandler "github.com/mirumyantsev/video_hosting/pkg/auth/handler"
	authrepo "github.com/mirumyantsev/video_hosting/pkg/auth/repository"
	authusecase "github.com/mirumyantsev/video_hosting/pkg/auth/usecase"
	"github.com/mirumyantsev/video_hosting/pkg/config"
	sconfig "github.com/mirumyantsev/video_hosting/pkg/config_stream"
	"github.com/mirumyantsev/video_hosting/pkg/download"
	downloadhandler "github.com/mirumyantsev/video_hosting/pkg/download/handler"
	downloadusecase "github.com/mirumyantsev/video_hosting/pkg/download/usecase"
	"github.com/mirumyantsev/video_hosting/pkg/logger"
	logrepo "github.com/mirumyantsev/video_hosting/pkg/logger/repository"
	logusecase "github.com/mirumyantsev/video_hosting/pkg/logger/usecase"
	"github.com/mirumyantsev/video_hosting/pkg/stream"
	streamhandler "github.com/mirumyantsev/video_hosting/pkg/stream/handler"
	streamrepo "github.com/mirumyantsev/video_hosting/pkg/stream/repository"
	streamusecase "github.com/mirumyantsev/video_hosting/pkg/stream/usecase"
	"github.com/mirumyantsev/video_hosting/pkg/user"
	userhandler "github.com/mirumyantsev/video_hosting/pkg/user/handler"
	userrepo "github.com/mirumyantsev/video_hosting/pkg/user/repository"
	userusecase "github.com/mirumyantsev/video_hosting/pkg/user/usecase"
)

type App struct {
	httpServer      *http.Server
	cfg             *config.Config
	scfg            *sconfig.Config
	userUseCase     user.UserUseCase
	authUseCase     auth.AuthUseCase
	sessUseCase     sess.SessUseCase
	logUseCase      logger.LogUseCase
	groupUseCase    group.GroupUseCase
	permUseCase     perm.PermUseCase
	infoUseCase     info.InfoUseCase
	videoUseCase    video.VideoUseCase
	StreamUC        stream.StreamUseCase
	downloadUseCase download.DownloadUseCase
}

func NewApp(cfg *config.Config, scfg *sconfig.Config) *App {
	userRepo := userrepo.NewUserRepository(cfg)
	authRepo := authrepo.NewAuthRepository(cfg)
	sessRepo := sessrepo.NewSessRepository(cfg)
	logRepo := logrepo.NewLogRepository(cfg)
	groupRepo := grouprepo.NewGroupRepository(cfg)
	permRepo := permrepo.NewPermRepository(cfg)
	infoRepo := inforepo.NewInfoRepository(cfg)
	videoRepo := videorepo.NewVideoRepository(cfg)
	streamRepo := streamrepo.NewStreamRepository(cfg)

	return &App{
		cfg:             cfg,
		scfg:            scfg,
		userUseCase:     userusecase.NewUserUseCase(cfg, userRepo),
		authUseCase:     authusecase.NewAuthUseCase(cfg, authRepo),
		sessUseCase:     sessusecase.NewSessUseCase(sessRepo, authRepo),
		logUseCase:      logusecase.NewLogUseCase(logRepo),
		groupUseCase:    groupusecase.NewGroupUseCase(groupRepo),
		permUseCase:     permusecase.NewPermUseCase(permRepo),
		infoUseCase:     infousecase.NewInfoUseCase(infoRepo),
		videoUseCase:    videousecase.NewVideoUseCase(videoRepo),
		StreamUC:        streamusecase.NewStreamUseCase(cfg, scfg, streamRepo),
		downloadUseCase: downloadusecase.NewDownloadUseCase(cfg),
	}
}

func (a *App) Run() error {
	// Debug mode
	if a.cfg.ServerDebugEnable {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init engine
	router := gin.New()

	// Init middleware
	router.Use(CORSMiddleware())

	// Check for web directory exists and register routes
	if _, err := os.Stat("./web"); !os.IsNotExist(err) {
		router.LoadHTMLGlob("./web/templates/*")
		streamhandler.RegisterTemplateHTTPEndpoints(router, a.cfg, a.scfg, a.StreamUC,
			a.userUseCase, a.logUseCase, a.authUseCase, a.sessUseCase)
	}

	router.StaticFS("/static", http.Dir("./web/static"))

	// Register routes
	authhandler.RegisterHTTPEndpoints(router, a.cfg, a.authUseCase, a.userUseCase,
		a.sessUseCase, a.logUseCase)
	userhandler.RegisterHTTPEndpoints(router, a.cfg, a.userUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase)
	grouphandler.RegisterHTTPEndpoints(router, a.cfg, a.groupUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)
	permhandler.RegisterHTTPEndpoints(router, a.cfg, a.permUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase, a.groupUseCase)
	infohandler.RegisterHTTPEndpoints(router, a.cfg, a.scfg, a.infoUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)
	videohandler.RegisterHTTPEndpoints(router, a.cfg, a.videoUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)
	streamhandler.RegisterStreamingHTTPEndpoints(router, a.cfg, a.scfg, a.StreamUC,
		a.userUseCase, a.logUseCase, a.authUseCase, a.sessUseCase)
	downloadhandler.RegisterHTTPEndpoints(router, a.cfg, a.downloadUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", a.cfg.ServerHost, a.cfg.ServerPort),
		Handler:        router,
		ReadTimeout:    time.Duration(a.cfg.ServerReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(a.cfg.ServerWriteTimeoutSeconds) * time.Second,
		MaxHeaderBytes: a.cfg.ServerMaxHeaderBytes,
	}

	// Server start
	var err error
	go func() {
		err = a.httpServer.ListenAndServe()
	}()
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot start server. Error: %s.", err.Error()))
	}
	a.cfg.ServerIP = getOutboundIP()
	logger.Print(msg.InfoServerStartedSuccessfullyAtLocalAddress(a.cfg.ServerIP, a.cfg.ServerPort))

	// Listening for interrupt signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Print(msg.InfoRecivedSignal(sig))
		done <- true
	}()
	<-done

	// Server shut down
	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return errors.New(fmt.Sprintf("Cannot shut down the server correctly. Error: %s.", err.Error()))
	}

	logger.Print(msg.InfoServerShutedDownCorrectly())

	return nil
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.Print(msg.WarningCannotGetLocalIP(err))
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
