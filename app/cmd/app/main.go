package main

import (
	"github.com/disbeliefff/ecommerce/internal/config"
	"log"
)

func main() {
	log.Println("[Config] initializing config")
	cfg := config.GetConfig()
}
