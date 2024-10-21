package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEBUG" env-default:"false"`
	Listen        struct {
		Type   string `env:"LISTEN_TYPE"    env-default:"port"`
		BindIP string `env:"LISTEN_BIND_IP" env-default:"0.0.0.0"`
		Port   string `env:"LISTEN_PORT"    env-default:"8081"`
	}
	AppConfig struct {
		LogLevel  string
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-required:"true"`
			Password string `env:"ADMIN_USER_PASSWORD" env-required:"true"`
		}
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("Collecting config")
		instance = &Config{}
		err := cleanenv.ReadEnv(instance)
		if err != nil {
			errText := "failed to collect config: %s"
			configDesc, _ := cleanenv.GetDescription(instance, &errText)
			log.Fatalf(configDesc, err)
		}
	})
	return instance
}
