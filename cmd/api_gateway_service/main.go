package main

import (
	"log"
	"os"

	"github.com/dinorain/kalobranded/config"
	"github.com/dinorain/kalobranded/internal/server"
	"github.com/dinorain/kalobranded/pkg/logger"
	"github.com/dinorain/kalobranded/pkg/postgres"
	"github.com/dinorain/kalobranded/pkg/redis"
	"github.com/dinorain/kalobranded/pkg/utils"
)

// @contact.name Dustin Jourdan
// @contact.url https://github.com/dinorain
// @contact.email djourdan555@gmail.com

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

func main() {
	log.Println("Starting kalobranded service")

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
	appLogger.Info("Postgresql connected")

	authServer := server.NewAppServer(appLogger, cfg, psqlDB, redisClient)
	appLogger.Fatal(authServer.Run())
}
