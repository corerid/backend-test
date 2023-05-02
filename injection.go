package main

import (
	"github.com/corerid/backend-test/config"
	"github.com/corerid/backend-test/db"
	"github.com/corerid/backend-test/handlers"
	httpX "github.com/corerid/backend-test/http"
	"github.com/corerid/backend-test/repositories"
	"github.com/corerid/backend-test/services"
	"net/http"
)

func injection(dbCon *db.Connection, config config.Config) handlers.HandlerI {
	repo := &repositories.Repository{DB: dbCon.RW}

	client := httpX.New(&http.Client{})
	service := &services.Service{
		RepositoryI: repo,
		Config:      config,
		Client:      client,
	}
	return &handlers.Handler{ServiceI: service}
}
