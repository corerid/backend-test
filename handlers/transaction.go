package handlers

import (
	"github.com/corerid/backend-test/services"
	"github.com/gin-gonic/gin"
	"math/big"
)

type getTransactionResponse struct {
	Transactions []services.Transaction `json:"transactions"`
	NextCursorID *int64                 `json:"next_cursor_id"`
}

type getTransactionQueryParams struct {
	StartBlock *big.Int `form:"start_block"`
	EndBlock   *big.Int `form:"end_block"`
	Address    string   `form:"address"`
	Limit      int      `form:"limit"`
	CursorID   int64    `form:"cursor_id"`
}

func (h *Handler) GetTransactionHandler(ctx *gin.Context) {
	getTransaction, err := parseRequestToService(ctx)
	if err != nil {
		ctx.JSON(400, err)
	}

	transaction, err := h.GetTransaction(getTransaction)
	if err != nil {
		ctx.JSON(400, err)
	}

	ctx.JSON(200, createResponse(transaction))

}

func parseRequestToService(ctx *gin.Context) (services.GetTransaction, error) {
	var queryParams getTransactionQueryParams
	err := ctx.BindQuery(&queryParams)
	if err != nil {
		return services.GetTransaction{}, err
	}

	return services.GetTransaction{
		Address:    queryParams.Address,
		StartBlock: queryParams.StartBlock,
		EndBlock:   queryParams.EndBlock,
		CursorID:   queryParams.CursorID,
		Limit:      queryParams.Limit,
	}, nil
}

func createResponse(transactions []services.Transaction) getTransactionResponse {
	var nextCursorId *int64
	if len(transactions) > 0 {
		nextCursorId = &transactions[len(transactions)-1].ID
	}

	return getTransactionResponse{
		Transactions: transactions,
		NextCursorID: nextCursorId,
	}
}
