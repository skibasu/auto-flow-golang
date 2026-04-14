package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skibasu/auto-flow-api/internal/appMiddleware"
	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/db"
	"github.com/skibasu/auto-flow-api/internal/handlers"
	"github.com/skibasu/auto-flow-api/internal/repository"
	"github.com/skibasu/auto-flow-api/internal/router"
	"github.com/skibasu/auto-flow-api/internal/services"
)

type Server struct {
	Config     config.Config
	DB         *pgxpool.Pool
	Handler    *handlers.Handler
	Router     *router.Router
	MiddleWare *appMiddleware.AppMiddleware
}

func NewServer(config config.Config) *Server {
	var s Server
	s.Config = config

	db, err := db.NewDB(s.Config.DBUrl)
	if err != nil {
		panic(err.Error())
	}
	s.DB = db
	s.MiddleWare = appMiddleware.NewAppMiddleware(config)
	repository := repository.NewRepository(s.DB)
	services := services.NewService(repository, config)

	s.Handler = handlers.NewHandler(services)
	s.Router = router.NewRouter()
	s.printBootstrapStatus()

	return &s

}

func (s *Server) RunServer() {
	server := &http.Server{
		Addr:              ":" + s.Config.AppPort,
		Handler:           s.Router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	s.printStartupBanner()
	err := server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("%s%s[error]%s Failed to listen the server: %v\n", ansiBold, ansiGreen, ansiReset, err)
	}
}
