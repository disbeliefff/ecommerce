package config

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
