package main

import (
	"github.com/disbeliefff/ecommerce/internal/app"
	"github.com/disbeliefff/ecommerce/internal/config"
	"github.com/disbeliefff/ecommerce/pkg/logging"
	"log"
)

func main() {
	log.Println("[Config] initializing config")
	cfg := config.GetConfig()

	log.Println("[Logger] initializing logger")
	logging.Init("debug")
	logger := logging.Init(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}

	log.Println("[Server] app starting")
	a.Run()
}
