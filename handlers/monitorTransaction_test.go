package handlers

import (
	mockServices "github.com/corerid/backend-test/services/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestHandler_MonitorBlockEthereumHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockServices.NewMockServiceI(ctrl)

	specifiedAddress := "0x123456789"
	mockService.EXPECT().MonitorBlockEthereum(specifiedAddress)

	h := &Handler{ServiceI: mockService}

	h.MonitorBlockEthereumHandler(specifiedAddress)
}
