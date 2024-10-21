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
	logger := logging.GetLogger()

	app, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}
}
