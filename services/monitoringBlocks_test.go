package services

import (
	"encoding/json"
	"errors"
	httpX "github.com/corerid/backend-test/http"
	mockHTTP "github.com/corerid/backend-test/http/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
)

func Test_getLatestBlock(t *testing.T) {
	tests := []struct {
		name         string
		prepPostBody getLatestBlockNumberRequest
		prepPostResp getLatestBlockNumberResponse
		postErr      error
		expected     *big.Int
		hasError     bool
	}{
		{
			name: "test_get_latest_block_success",
			prepPostBody: getLatestBlockNumberRequest{
				Jsonrpc: "2.0",
				Method:  "eth_blockNumber",
				Id:      0,
			},
			prepPostResp: getLatestBlockNumberResponse{
				Jsonrpc: "2.0",
				Result:  "0x1058462",
				Id:      0,
			},
			postErr:  nil,
			expected: big.NewInt(17138786),
			hasError: false,
		},
		{
			name: "test_get_latest_block_fail_when_send_error",
			prepPostBody: getLatestBlockNumberRequest{
				Jsonrpc: "2.0",
				Method:  "eth_blockNumber",
				Id:      0,
			},
			postErr:  errors.New("cannot send request"),
			hasError: true,
		},
		{
			name: "test_get_latest_block_fail_when_convert_hex_to_big_int_error",
			prepPostBody: getLatestBlockNumberRequest{
				Jsonrpc: "2.0",
				Method:  "eth_blockNumber",
				Id:      0,
			},
			prepPostResp: getLatestBlockNumberResponse{},
			postErr:      nil,
			hasError:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httptest := mockHTTP.NewMockClient(ctrl)

			sendReq, err := json.Marshal(tt.prepPostBody)
			assert.NoError(t, err)

			postResp, err := json.Marshal(tt.prepPostResp)
			assert.NoError(t, err)

			respReadCloser := io.NopCloser(strings.NewReader(string(postResp)))
			postRespBody := &http.Response{
				Body: respReadCloser,
			}

			mockURL := "http://test"
			mockReq := httpX.NewRequest(sendReq, "application/json")
			httptest.EXPECT().Send(http.MethodPost, mockURL, mockReq).Return(postRespBody, tt.postErr)

			actual, err := getLatestBlock(mockURL, httptest)

			if tt.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_getBlockByNumber(t *testing.T) {

	tests := []struct {
		name         string
		blockNo      *big.Int
		prepPostBody getBlockByNumberRequest
		expectedResp getBlockByNumberResponse
		postErr      error
		hasError     bool
	}{
		{
			name:    "test_get_block_by_number_success",
			blockNo: big.NewInt(17138786),
			prepPostBody: getBlockByNumberRequest{
				Jsonrpc: "2.0",
				Method:  "eth_getBlockByNumber",
				Params: []interface{}{
					"0x1058462",
					true,
				},
				Id: 0,
			},
			expectedResp: getBlockByNumberResponse{
				Jsonrpc: "2.0",
				ID:      0,
			},
			postErr:  nil,
			hasError: false,
		},
		{
			name:    "test_get_block_by_number_fail_when_send_request_is_error",
			blockNo: big.NewInt(17138786),
			prepPostBody: getBlockByNumberRequest{
				Jsonrpc: "2.0",
				Method:  "eth_getBlockByNumber",
				Params: []interface{}{
					"0x1058462",
					true,
				},
				Id: 0,
			},
			expectedResp: getBlockByNumberResponse{
				Jsonrpc: "2.0",
				ID:      0,
			},
			postErr:  errors.New("cannot send request"),
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httptest := mockHTTP.NewMockClient(ctrl)

			sendReq, err := json.Marshal(tt.prepPostBody)
			assert.NoError(t, err)

			postResp, err := json.Marshal(tt.expectedResp)
			assert.NoError(t, err)

			respReadCloser := io.NopCloser(strings.NewReader(string(postResp)))
			postRespBody := &http.Response{
				Body: respReadCloser,
			}

			mockURL := "http://test"
			mockReq := httpX.NewRequest(sendReq, "application/json")
			httptest.EXPECT().Send(http.MethodPost, mockURL, mockReq).Return(postRespBody, tt.postErr)

			actual, err := getBlockByNumber(mockURL, httptest, tt.blockNo)

			if tt.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp, actual)
		})
	}
}
