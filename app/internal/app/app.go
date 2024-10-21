package app

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/disbeliefff/ecommerce/docs"
	"github.com/disbeliefff/ecommerce/internal/config"
	"github.com/disbeliefff/ecommerce/pkg/logging"
	"github.com/disbeliefff/ecommerce/pkg/metric"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net"
	"net/http"
	"time"
)

type App struct {
	cfg        *config.Config
	lg         *logging.Logger
	router     *chi.Mux
	httpServer *http.Server
}

func NewApp(cfg *config.Config, lg *logging.Logger) (*App, error) {
	router := chi.NewRouter()

	lg.Println("[swagger] initializing swagger")
	router.Get("/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently).ServeHTTP)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("[metric] initializing heartbeat metric")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	app := &App{
		cfg:    cfg,
		lg:     lg,
		router: router,
	}
	return app, nil
}

func (a *App) Run() {
	a.lg.Println("[server]Starting server")
	a.StartHTTP()
}

func (a *App) StartHTTP() error {

	logger := a.lg.LWithFields(map[string]any{
		"IP":   a.cfg.Listen.BindIP,
		"Port": a.cfg.Listen.Port,
	})
	logger.Info("[server]HTTP Server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
	if err != nil {
		logger.LWithField("error", err).Fatal("failed to create listener")
	}

	logger.Info("[server]CORS initializing")
	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:8081"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Authorization", "Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Access-Token", "Refresh-Token", "Location", "Authorization", "Content-Disposition"},
		Debug:              false,
	})

	a.httpServer = &http.Server{
		Handler:      c.Handler(a.router),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("[server]Starting HTTP Server")
	if err = a.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.LWithField("error", err).Fatal("failed to start server")
		return err
	}

	if err := a.ShutdownHTTP(); err != nil {
		logger.LWithField("error", err).Fatal("failed to shutdown server")
		return err
	}

	logger.Info("[server]Shutting down HTTP server")
	return nil
}

func (a *App) ShutdownHTTP() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.lg.Info("[server]Gracefully shutting down the HTTP server")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.lg.LWithField("error", err).Error("failed to gracefully shutdown the server")
		return err
	}

	a.lg.Info("[server]HTTP server successfully shut down")
	return nil
}
