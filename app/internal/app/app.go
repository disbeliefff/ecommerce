package app

import (
	"github.com/disbeliefff/ecommerce/internal/config"
	"github.com/disbeliefff/ecommerce/pkg/logging"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type App struct {
	cfg *config.Config
	lg  *logging.Logger
}

func NewApp(cfg *config.Config, lg *logging.Logger) (*App, error) {
	r := chi.NewRouter()

	lg.Println("[swagger] initializing swagger")
	r.Get("/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	app := &App{
		cfg: cfg,
		lg:  lg,
	}
	return app, nil
}
