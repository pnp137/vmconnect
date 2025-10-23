package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/logger"
	"linksupply.io/vmconnect/system"

	"github.com/gofiber/fiber/v2"
)

type APIServer struct {
	appConfig  *system.Config
	logger     *logger.AppLogger
	dataSource *system.DataSource
	app        *fiber.App
}

func NewServer(appConfig *system.Config) *APIServer {
	logger := logger.NewLogger()
	dataSource := system.NewDataSource(&appConfig.Db, appConfig.LogLevel == "DEBUG")
	app := fiber.New(fiber.Config{
		CaseSensitive:         true,
		ServerHeader:          "esper/apps-api",
		DisableStartupMessage: true,
		ErrorHandler:          response.DefaultErrorHandler,
		BodyLimit:             appConfig.BodyUploadLimit,
	})
	return &APIServer{
		appConfig:  appConfig,
		logger:     logger,
		dataSource: dataSource,
		app:        app,
	}
}

func (server *APIServer) StartServer() {
	log := server.logger
	listenPort := fmt.Sprintf(":%d", server.appConfig.ListenPort)
	SetupRoutes(server)
	log.Info(fmt.Sprintf("Starting webserver for vmconnect-api provider on port %s", listenPort))
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := server.app.Listen(listenPort)
			if err != nil {
				log.Panic(err, fmt.Sprintf("Could not start API server: %s", err))
			}
		}
	}()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Info("Terminating api server: context cancelled")
		server.app.Shutdown()
	case <-sigterm:
		log.Info("Terminating api server: via signal")
		server.app.Shutdown()
	}
	cancel()
}
