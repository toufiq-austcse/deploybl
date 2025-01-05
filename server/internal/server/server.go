package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/config"
)

type Server struct {
	GinEngine *gin.Engine
	Server    *http.Server
}

func NewServer() *Server {
	r := gin.Default()
	server := &http.Server{
		Addr:    ":" + config.AppConfig.PORT,
		Handler: r,
	}

	return &Server{r, server}
}

func (s *Server) Run() error {
	return s.GinEngine.Run(":" + config.AppConfig.PORT)
}

//func (s *Server) Stop() error {
//	return s.GinEngine.SHU
//}
