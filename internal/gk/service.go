package gk

import (
	"errors"
)

var (
	// ErrServiceNotExist error when getting not registered service
	ErrServiceNotExist = errors.New("service does not exist")
	// ErrServiceAlreadyExists error when setting already registered service
	ErrServiceAlreadyExists = errors.New("service already exists")
)

// Service interface for services
type Service interface {
	Init() error
	Up() error
	Down()
}

// ServiceContainer interface for contained services
type ServiceContainer interface {
	Get(string) (Service, error)
	Set(string, Service) error
}
