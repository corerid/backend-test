package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/corerid/backend-test/utils"
	"math/big"
	"strconv"
	"strings"
)

type Transaction struct {
	ID              int64
	TransactionHash string
	Address         string
	BlockNumber     *big.Int
	Data            json.RawMessage
}

func (repo Repository) CreateTransactionRepo(transaction Transaction) error {
	createTransactionSQLCommand := `INSERT INTO transaction (transaction_hash, address, block_number, "data")
							VALUES($1, $2, $3, $4)`

	_, err := repo.DB.Exec(createTransactionSQLCommand, transaction.TransactionHash, transaction.Address, transaction.BlockNumber.String(), transaction.Data)
	if err != nil {
		return err
	}

	return nil
}

func (repo Repository) GetTransactionRepo(address string, startBlock *big.Int, endBlock *big.Int, cursorID int64, limit int) ([]Transaction, error) {
	var transactions []Transaction

	getTransactionSQLCommand, queryArgs := createGetTransactionSQLStmt(address, startBlock, endBlock, cursorID, limit)

	query, err := repo.DB.Prepare(getTransactionSQLCommand)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	rows, err := query.Query(queryArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var transaction Transaction

		var blockNumberBytes []byte

		err = rows.Scan(&transaction.ID, &transaction.TransactionHash, &transaction.Address,
			&blockNumberBytes, &transaction.Data)
		if err != nil {
			return nil, err
		}

		blockNumber := new(big.Int)
		blockNumber.SetString(string(blockNumberBytes), 10)

		transaction.BlockNumber = blockNumber

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func createGetTransactionSQLStmt(address string, startBlock *big.Int, endBlock *big.Int, cursorID int64, limit int) (string, []interface{}) {
	getTransactionSQLCommand := `SELECT id, transaction_hash, address, block_number, data 
								 FROM transaction %s
								 ORDER BY id ASC
								 %s`

	var queryArgs []interface{}

	var whereClauseList []string
	if address != "" {
		whereClauseList = append(whereClauseList, fmt.Sprintf("address = $%s", strconv.Itoa(len(queryArgs)+1)))
		queryArgs = append(queryArgs, address)
	}

	startBlockCBI := utils.CustomBigInt{Int: startBlock}
	if !startBlockCBI.IsEmptyOrNil() {
		whereClauseList = append(whereClauseList, fmt.Sprintf("block_number >= $%s", strconv.Itoa(len(queryArgs)+1)))
		queryArgs = append(queryArgs, startBlock.String())
	}

	endBlockCBI := utils.CustomBigInt{Int: endBlock}
	if !endBlockCBI.IsEmptyOrNil() {
		whereClauseList = append(whereClauseList, fmt.Sprintf("block_number <= $%s", strconv.Itoa(len(queryArgs)+1)))
		queryArgs = append(queryArgs, endBlock.String())
	}

	if cursorID > 0 {
		whereClauseList = append(whereClauseList, fmt.Sprintf("id > $%s", strconv.Itoa(len(queryArgs)+1)))
		queryArgs = append(queryArgs, cursorID)
	}

	if len(whereClauseList) > 0 {
		whereClauseList[0] = "WHERE " + whereClauseList[0]
	}
	whereClause := strings.Join(whereClauseList, "\nAND ")

	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("\nLIMIT $%s", strconv.Itoa(len(queryArgs)+1))
		queryArgs = append(queryArgs, limit)
	}

	getTransactionSQLCommand = fmt.Sprintf(getTransactionSQLCommand, whereClause, limitClause)

	return getTransactionSQLCommand, queryArgs
}
