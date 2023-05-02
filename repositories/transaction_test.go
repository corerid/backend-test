package repositories

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"regexp"
	"testing"
)

func TestRepository_GetTransactionRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	address := "0x123"
	startBlock := big.NewInt(10)
	endBlock := big.NewInt(20)
	cursorID := int64(1)
	limit := 10
	expectedQuery := `SELECT id, transaction_hash, address, block_number, data 
					  FROM transaction WHERE address = $1
					  AND block_number >= $2
					  AND block_number <= $3
					  AND id > $4
					  ORDER BY id ASC			 
					  LIMIT $5`
	expectedArgs := []driver.Value{address, startBlock.String(), endBlock.String(), cursorID, limit}

	expectedResult := []Transaction{
		Transaction{
			ID:              1,
			TransactionHash: "0x1111",
			Address:         "0x9999",
			BlockNumber:     big.NewInt(1),
			Data:            json.RawMessage("test 1"),
		},
		Transaction{
			ID:              2,
			TransactionHash: "0x2222",
			Address:         "0x9999",
			BlockNumber:     big.NewInt(2),
			Data:            json.RawMessage("test 2"),
		},
	}

	mock.ExpectPrepare(regexp.QuoteMeta(expectedQuery))

	rows := sqlmock.NewRows([]string{"id", "transaction_hash", "address", "block_number", "data"}).
		AddRow(expectedResult[0].ID, expectedResult[0].TransactionHash, expectedResult[0].Address, expectedResult[0].BlockNumber.String(), expectedResult[0].Data).
		AddRow(expectedResult[1].ID, expectedResult[1].TransactionHash, expectedResult[1].Address, expectedResult[1].BlockNumber.String(), expectedResult[1].Data)

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(expectedArgs...).WillReturnRows(rows)

	repo := Repository{DB: db}

	result, err := repo.GetTransactionRepo(address, startBlock, endBlock, cursorID, limit)
	assert.NoError(t, err)

	assert.Equal(t, len(expectedResult), len(result))

	for i, r := range result {
		assert.Equal(t, expectedResult[i].ID, r.ID)
		assert.Equal(t, expectedResult[i].TransactionHash, r.TransactionHash)
		assert.Equal(t, expectedResult[i].Address, r.Address)
		assert.Equal(t, 0, r.BlockNumber.Cmp(expectedResult[i].BlockNumber))
		assert.Equal(t, string(expectedResult[i].Data), string(r.Data))
	}

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRepository_CreateTransactionRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	transaction := Transaction{
		TransactionHash: "0x1111",
		Address:         "0x9999",
		BlockNumber:     big.NewInt(1),
		Data:            []byte("test"),
	}

	expectedQuery := `INSERT INTO transaction (transaction_hash, address, block_number, "data")
					  VALUES($1, $2, $3, $4)`

	expectedArgs := []driver.Value{transaction.TransactionHash, transaction.Address, transaction.BlockNumber.String(), transaction.Data}

	mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WithArgs(expectedArgs...).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := Repository{DB: db}
	err = repo.CreateTransactionRepo(transaction)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
