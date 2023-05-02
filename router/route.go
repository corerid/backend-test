package router

import (
	"github.com/corerid/backend-test/handlers"
	"github.com/gin-gonic/gin"
)

func AddRoute(h handlers.HandlerI) *gin.Engine {
	r := gin.Default()

	r.GET("/transaction", h.GetTransactionHandler)

	return r
}
