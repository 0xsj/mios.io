package api

import (
	"github.com/0xsj/gin-sqlc/config"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	router *gin.Engine
	store db.Querier
	log log.Logger
}

func NewServer(config config.Config, store db.Querier, log log.Logger) *Server {
	server := &Server{
		config: config,
		router: gin.Default(),
		store:  store,
		log:    log,
	}
	return server
}
