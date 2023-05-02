package services

import (
	"github.com/corerid/backend-test/config"
	"github.com/corerid/backend-test/http"
	"github.com/corerid/backend-test/repositories"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go

type ServiceI interface {
	GetTransaction(getTransaction GetTransaction) ([]Transaction, error)
	MonitorBlockEthereum(specifiedAddress string)
}

type Service struct {
	repositories.RepositoryI
	config.Config
	http.Client
}
