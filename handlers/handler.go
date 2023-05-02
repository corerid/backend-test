package handlers

import (
	"github.com/corerid/backend-test/services"
	"github.com/gin-gonic/gin"
)

type HandlerI interface {
	GetTransactionHandler(ctx *gin.Context)
	MonitorBlockEthereumHandler(specifiedAddress string)
}

type Handler struct {
	services.ServiceI
}
