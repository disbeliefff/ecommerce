package app

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/disbeliefff/ecommerce/docs"
	"github.com/disbeliefff/ecommerce/internal/config"
	"github.com/disbeliefff/ecommerce/pkg/client/postgresql"
	"github.com/disbeliefff/ecommerce/pkg/logging"
	"github.com/disbeliefff/ecommerce/pkg/metric"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *chi.Mux
	httpServer *http.Server
	pgClient   *pgxpool.Pool
}

func NewApp(cfg *config.Config, lg *logging.Logger) (*App, error) {
	router := chi.NewRouter()

	lg.Println("[swagger] initializing swagger")
	router.Get("/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently).ServeHTTP)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("[metric] initializing heartbeat metric")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	pgConfig := postgresql.NewPgConfig(cfg.PostgreSQL.Username, cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Database)

	pgClient, err := postgresql.NewClient(context.Background(), 5, 5, pgConfig)

	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		cfg:      cfg,
		logger:   lg,
		router:   router,
		pgClient: pgClient,
	}
	return app, nil
}

func (a *App) Run() {
	a.logger.Println("[server]Starting server")
	a.startHTTP()
}

func (a *App) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("application completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}
}
