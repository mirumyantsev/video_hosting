package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikerumy/vhosting/internal/auth"
	authhandler "github.com/mikerumy/vhosting/internal/auth/handler"
	authrepo "github.com/mikerumy/vhosting/internal/auth/repository"
	authusecase "github.com/mikerumy/vhosting/internal/auth/usecase"
	logging "github.com/mikerumy/vhosting/internal/logging"
	loggingrepo "github.com/mikerumy/vhosting/internal/logging/repository"
	loggingusecase "github.com/mikerumy/vhosting/internal/logging/usecase"
	"github.com/mikerumy/vhosting/internal/models"
	"github.com/mikerumy/vhosting/internal/user"
	userhandler "github.com/mikerumy/vhosting/internal/user/handler"
	userrepo "github.com/mikerumy/vhosting/internal/user/repository"
	userusecase "github.com/mikerumy/vhosting/internal/user/usecase"
)

type App struct {
	httpServer     *http.Server
	cfg            models.Config
	userUseCase    user.UserUseCase
	authUseCase    auth.AuthUseCase
	loggingUseCase logging.LoggingUseCase
}

func NewApp(cfg models.Config) *App {
	userRepo := userrepo.NewUserRepository(cfg)
	authRepo := authrepo.NewAuthRepository(cfg)
	loggingRepo := loggingrepo.NewLoggingRepository(cfg)

	return &App{
		cfg:            cfg,
		userUseCase:    userusecase.NewUserUseCase(cfg, userRepo),
		authUseCase:    authusecase.NewAuthUseCase(cfg, authRepo),
		loggingUseCase: loggingusecase.NewLoggingUseCase(cfg, loggingRepo),
	}
}

func (a *App) Run() error {
	// Server mode
	if a.cfg.ServerDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init Gin handler
	router := gin.New()

	// Register web routes
	authhandler.RegisterHTTPEndpoints(router, a.authUseCase, a.userUseCase, a.loggingUseCase)
	userhandler.RegisterHTTPEndpoints(router, a.userUseCase)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + a.cfg.ServerPort,
		Handler:        router,
		ReadTimeout:    time.Duration(a.cfg.ServerReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(a.cfg.ServerWriteTimeoutSeconds) * time.Second,
		MaxHeaderBytes: a.cfg.ServerMaxHeaderBytes,
	}

	var err error
	notStarted := false

	// Server start
	go func() {
		err = a.httpServer.ListenAndServe()
		if err != nil {
			notStarted = true
		}
	}()
	time.Sleep(50 * time.Millisecond)
	if notStarted {
		return errors.New(fmt.Sprintf("Cannot start server. Error: %s.", err.Error()))
	}
	InfoServerWasSuccessfullyStartedAtLocalIP(getOutboundIP().String(), a.cfg.ServerPort)

	// Server shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	time.Sleep(2 * time.Second)

	ctx, shutdown := context.WithTimeout(context.Background(), 1700*time.Millisecond)
	defer shutdown()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return errors.New(fmt.Sprintf("Cannot shut down the server correctly. Error: %s.", err.Error()))
	}

	InfoServerWasGracefullyShutDown()
	return nil
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		WarningCannotGetLocalIP(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
