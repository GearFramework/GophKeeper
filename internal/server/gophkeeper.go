package server

import (
	"github.com/GearFramework/GophKeeper/internal/config"
	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/GearFramework/GophKeeper/internal/pkg/entity"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/rest"
	"github.com/GearFramework/GophKeeper/internal/pkg/storage/db"
	"github.com/GearFramework/GophKeeper/internal/pkg/user"
	"github.com/GearFramework/GophKeeper/internal/server/v1"
)

type servicesMap = map[string]gk.Service

// GkServer struct of app server
type GkServer struct {
	conf     *config.Server
	services servicesMap
}

// NewGkServer retuning GopKeeper app server
func NewGkServer(conf *config.Server) *GkServer {
	return &GkServer{
		conf:     conf,
		services: servicesMap{},
	}
}

// Init GophKeeper
func (gks *GkServer) Init() error {
	logger.Log.Info("initializing GophKeeper")
	rs := rest.NewServer(&rest.Config{
		Addr: gks.conf.Addr,
	}, v1.NewRestAPI(gks.conf))
	if err := gks.Set("rest", rs); err != nil {
		return err
	}
	store := db.NewStorage(gks.conf.StorageDSN)
	if err := gks.Set("store", store); err != nil {
		return err
	}
	return gks.initRepositories()
}

func (gks *GkServer) initRepositories() error {
	store, err := gks.Get("store")
	if err != nil {
		return err
	}
	user.NewRepository(store.(gk.DBStorable))
	entity.NewRepository(store.(gk.DBStorable))
	return nil
}

// Run server
func (gks *GkServer) Run() error {
	logger.Log.Info("GophKeeper run")
	errChan := make(chan error)
	for name, service := range gks.services {
		go func() {
			logger.Log.Infof("up service %s", name)
			errChan <- service.Up()
		}()
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}

// Stop shutdown handler
func (gks *GkServer) Stop() {
	for _, service := range gks.services {
		service.Down()
	}
}

// Get service by name
func (gks *GkServer) Get(serviceName string) (gk.Service, error) {
	if s, ok := gks.services[serviceName]; ok {
		return s, nil
	}
	return nil, gk.ErrServiceNotExist
}

// Set service
func (gks *GkServer) Set(serviceName string, s gk.Service) error {
	if _, ok := gks.services[serviceName]; ok {
		return gk.ErrServiceAlreadyExists
	}
	if err := s.Init(); err != nil {
		return err
	}
	gks.services[serviceName] = s
	return nil
}

// GetConfig returning app config
func (gks *GkServer) GetConfig() *config.Server {
	return gks.conf
}
