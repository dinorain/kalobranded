package main

import (
	"log"
	"os"

	"github.com/dinorain/checkoutaja/config"
	"github.com/dinorain/checkoutaja/internal/server"
	"github.com/dinorain/checkoutaja/pkg/logger"
	"github.com/dinorain/checkoutaja/pkg/postgres"
	"github.com/dinorain/checkoutaja/pkg/redis"
	"github.com/dinorain/checkoutaja/pkg/utils"
)

// @contact.name Dustin Jourdan
// @contact.url https://github.com/dinorain
// @contact.email djourdan555@gmail.com

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

func main() {
	log.Println("Starting checkoutaja service")

	configPath := utils.GetConfigPath(os.Getenv("CONFIG"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewAppLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Server.AppVersion)

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	defer psqlDB.Close()

	authServer := server.NewAppServer(appLogger, cfg, psqlDB, redisClient)
	appLogger.Fatal(authServer.Run())
}
