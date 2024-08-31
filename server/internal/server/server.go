package server

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/config"
)

type Server struct {
	GinEngine *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()
	return &Server{r}
}

func (s *Server) Run() error {
	return s.GinEngine.Run(":" + config.AppConfig.PORT)
}

//func (s *Server) Stop() error {
//	return s.GinEngine.SHU
//}
