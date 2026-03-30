package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/router"
)

type Server struct {
	Config  config.Config
	DB      *pgxpool.Pool
	Handler any
	Router  *router.Router
}
