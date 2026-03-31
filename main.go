package main

import (
	"github.com/skibasu/auto-flow-api/internal/config"
	"github.com/skibasu/auto-flow-api/internal/server"
)

func main() {
	cfg := config.NewConfig()
	srv := server.NewServer(cfg)

	defer srv.DB.Close()

	var router = srv.Router
	var handlers = srv.Handler
	var middleware = srv.MiddleWare

	//Initialize routes
	router.InitializeRoutes(handlers, middleware)
	//Run server
	srv.RunServer()

}
