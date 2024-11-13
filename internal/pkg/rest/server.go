package rest

import (
	"net/http"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// MiddlewareFunc type of middleware function
type MiddlewareFunc func() gin.HandlerFunc

// Server struct of REST server
type Server struct {
	http *http.Server
	conf *Config
	api  gk.API
}

// NewServer returning instance of REST server
func NewServer(conf *Config, api gk.API) *Server {
	return &Server{
		conf: conf,
		api:  api,
	}
}

// WithAPI append API server
func (s *Server) WithAPI(api gk.API) {
	s.api = api
}

// Init REST server
func (s *Server) Init() error {
	logger.Log.Info("initializing REST server")
	s.api.Init()
	return nil
}

// Up REST server
func (s *Server) Up() error {
	return s.upTLS()
}

var tlsCertFile = ".cert/certbundle.pem"
var tlsKeyFile = ".cert/server.key"

func (s *Server) upTLS() error {
	s.http = &http.Server{
		Addr:    s.conf.Addr,
		Handler: s.api.GetRouter().(*gin.Engine),
	}
	logger.Log.Infof("Start HTTPS server at the %s\n", s.conf.Addr)
	err := s.api.GetRouter().(*gin.Engine).RunTLS(s.conf.Addr, tlsCertFile, tlsKeyFile)
	if err != nil {
		logger.Log.Infof("Failed to Listen and Serve HTTPS: %v\n", err)
		return err
	}
	return nil
}

// Down REST server
func (s *Server) Down() {
	logger.Log.Info("Shutdown REST server")
}
