package handlers

import (
	"encoding/json"
	"errors"
	"github.com/corerid/backend-test/services"
	mockServices "github.com/corerid/backend-test/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandler_GetTransactionHandler(t *testing.T) {

	type args struct {
		queryParams    map[string]string
		getTransaction services.GetTransaction
	}
	tests := []struct {
		name                 string
		args                 args
		prepTransaction      []services.Transaction
		serviceError         error
		expectedStatusCode   int
		expectedResponseBody getTransactionResponse
		hasError             bool
	}{
		{
			name: "test_transaction_handler_success",
			args: args{
				queryParams: map[string]string{
					"start_block": "1",
					"end_block":   "4",
					"limit":       "2",
					"cursor_id":   "1",
					"address":     "0x12345678",
				},
				getTransaction: services.GetTransaction{
					Address:    "0x12345678",
					StartBlock: big.NewInt(1),
					EndBlock:   big.NewInt(4),
					CursorID:   1,
					Limit:      2,
				},
			},
			prepTransaction: []services.Transaction{
				{
					ID:              1,
					TransactionHash: "0x1111",
					Address:         "0x9999",
					BlockNumber:     big.NewInt(1),
					Data:            json.RawMessage("{ \"data\": \"test 1\" }"),
				},
				{
					ID:              2,
					TransactionHash: "0x2222",
					Address:         "0x9999",
					BlockNumber:     big.NewInt(2),
					Data:            json.RawMessage("{ \"data\": \"test 2\" }"),
				},
			},
			serviceError:       nil,
			expectedStatusCode: http.StatusOK,
			hasError:           false,
		},
		{
			name: "test_transaction_handler_fail_cannot_convert_limit",
			args: args{
				queryParams: map[string]string{
					"start_block": "1",
					"end_block":   "4",
					"limit":       "a",
					"cursor_id":   "1",
					"address":     "0x12345678",
				},
			},
			prepTransaction:    []services.Transaction{},
			expectedStatusCode: http.StatusBadRequest,
			hasError:           false,
		},
		{
			name: "test_transaction_handler_fail_service_error",
			args: args{
				queryParams: map[string]string{
					"start_block": "1",
					"end_block":   "4",
					"limit":       "2",
					"cursor_id":   "1",
					"address":     "0x12345678",
				},
				getTransaction: services.GetTransaction{
					Address:    "0x12345678",
					StartBlock: big.NewInt(1),
					EndBlock:   big.NewInt(4),
					CursorID:   1,
					Limit:      2,
				},
			},
			prepTransaction:    []services.Transaction{},
			serviceError:       errors.New("service error"),
			expectedStatusCode: http.StatusBadRequest,
			hasError:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{
				Header: make(http.Header),
				URL:    &url.URL{},
			}

			u := url.Values{}
			for key, val := range tt.args.queryParams {
				u.Add(key, val)
			}

			ctx.Request.Method = "GET"
			ctx.Request.Header.Set("Content-Type", "application/json")

			// set query params
			ctx.Request.URL.RawQuery = u.Encode()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mockServices.NewMockServiceI(ctrl)
			mockService.EXPECT().GetTransaction(tt.args.getTransaction).Return(tt.prepTransaction, tt.serviceError)

			h := &Handler{
				ServiceI: mockService,
			}

			h.GetTransactionHandler(ctx)

			assert.EqualValues(t, tt.expectedStatusCode, w.Code)

			// verify response body
			if len(tt.prepTransaction) > 0 {
				expectTransaction := getTransactionResponse{
					Transactions: tt.prepTransaction,
					NextCursorID: &tt.prepTransaction[len(tt.prepTransaction)-1].ID,
				}
				expectedResponseBody, _ := json.Marshal(expectTransaction)
				assert.Equal(t, expectedResponseBody, w.Body.Bytes())
			}
		})
	}
}
