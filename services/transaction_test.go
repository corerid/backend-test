package services

import (
	"errors"
	"github.com/corerid/backend-test/repositories"
	mockRepositories "github.com/corerid/backend-test/repositories/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestService_GetTransaction(t *testing.T) {

	type args struct {
		getTransaction GetTransaction
	}
	tests := []struct {
		name                string
		args                args
		expectedTransaction []repositories.Transaction
		expectedError       error
		hasError            bool
	}{
		{
			name: "test_get_transaction_success",
			args: args{
				getTransaction: GetTransaction{
					Address:    "0x9999",
					StartBlock: big.NewInt(1),
					EndBlock:   big.NewInt(4),
					CursorID:   1,
					Limit:      2,
				},
			},
			expectedTransaction: []repositories.Transaction{
				{
					TransactionHash: "0x1111",
					Address:         "0x8888",
					BlockNumber:     big.NewInt(2),
					Data:            []byte("test"),
				},
				{
					TransactionHash: "0x1111",
					Address:         "0x9999",
					BlockNumber:     big.NewInt(3),
					Data:            []byte("test 2"),
				},
			},
			hasError: false,
		},
		{
			name: "test_get_transaction_fail_get_transaction_repo_error",
			args: args{
				getTransaction: GetTransaction{
					Address:    "0x9999",
					StartBlock: big.NewInt(1),
					EndBlock:   big.NewInt(4),
					CursorID:   1,
					Limit:      2,
				},
			},
			expectedTransaction: []repositories.Transaction{},
			expectedError:       errors.New("cannot get transaction repo"),
			hasError:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepositories.NewMockRepositoryI(ctrl)

			mockRepo.EXPECT().
				GetTransactionRepo(tt.args.getTransaction.Address, tt.args.getTransaction.StartBlock,
					tt.args.getTransaction.EndBlock, tt.args.getTransaction.CursorID, tt.args.getTransaction.Limit).
				Return(tt.expectedTransaction, tt.expectedError)

			s := Service{
				RepositoryI: mockRepo,
			}

			actualTransactions, err := s.GetTransaction(tt.args.getTransaction)

			assert.Equal(t, len(tt.expectedTransaction), len(actualTransactions))
			if tt.hasError {
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)
			for i, expectedTransaction := range tt.expectedTransaction {
				assert.Equal(t, expectedTransaction.TransactionHash, actualTransactions[i].TransactionHash)
				assert.Equal(t, expectedTransaction.Address, actualTransactions[i].Address)
				assert.Equal(t, expectedTransaction.BlockNumber, actualTransactions[i].BlockNumber)
				assert.Equal(t, expectedTransaction.Data, actualTransactions[i].Data)
			}
		})
	}
}
