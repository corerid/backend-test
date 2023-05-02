package repositories

import (
	"database/sql"
	"math/big"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

type RepositoryI interface {
	GetTransactionRepo(address string, startBlock *big.Int, endBlock *big.Int, cursorID int64, limit int) ([]Transaction, error)
	CreateTransactionRepo(transaction Transaction) error
}

type Repository struct {
	DB *sql.DB
}
