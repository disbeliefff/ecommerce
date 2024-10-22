package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEBUG" env-default:"false"`
	Listen        struct {
		Type       string `env:"LISTEN_TYPE"    env-default:"port"`
		BindIP     string `env:"LISTEN_BIND_IP" env-default:"0.0.0.0"`
		Port       string `env:"LISTEN_PORT"    env-default:"8080"`
		SocketFile string `env:"SOCKET_FILE"    env-default:"./app.sock"`
	}
	AppConfig struct {
		LogLevel  string
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-default:"admin"`
			Password string `env:"ADMIN_USER_PASSWORD" env-default:"dev"`
		}
	}
	PostgreSQL struct {
		Username string `env:"PSQL_USERNAME" env-required:"true"`
		Password string `env:"PSQL_PASSWORD" env-required:"true"`
		Host     string `env:"PSQL_HOST" env-required:"true"`
		Port     string `env:"PSQL_PORT" env-required:"true"`
		Database string `env:"PSQL_DATABASE" env-required:"true"`
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

func (c *Config) GetListenAddress() string {
	return fmt.Sprintf("%s:%s", c.Listen.BindIP, c.Listen.Port)
}
