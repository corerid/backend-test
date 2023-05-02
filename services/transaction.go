package services

import (
	"encoding/json"
	"github.com/corerid/backend-test/repositories"
	"math/big"
)

type GetTransaction struct {
	Address    string
	StartBlock *big.Int
	EndBlock   *big.Int
	CursorID   int64
	Limit      int
}

type Transaction struct {
	ID              int64
	TransactionHash string
	Address         string
	BlockNumber     *big.Int
	Data            json.RawMessage
}

func (s Service) GetTransaction(getTransaction GetTransaction) ([]Transaction, error) {
	transaction, err := s.GetTransactionRepo(getTransaction.Address, getTransaction.StartBlock, getTransaction.EndBlock, getTransaction.CursorID, getTransaction.Limit)
	if err != nil {
		return nil, err
	}

	return parseTransactionRepoToTransaction(transaction), nil
}

func parseTransactionRepoToTransaction(transactionsRepo []repositories.Transaction) []Transaction {
	var transactions []Transaction
	for _, transactionRepo := range transactionsRepo {
		transaction := Transaction{
			ID:              transactionRepo.ID,
			TransactionHash: transactionRepo.TransactionHash,
			Address:         transactionRepo.Address,
			BlockNumber:     transactionRepo.BlockNumber,
			Data:            transactionRepo.Data,
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}
