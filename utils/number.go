package utils

import (
	"errors"
	"fmt"
	"math/big"
)

type CustomBigInt struct {
	*big.Int
}

func (cbi *CustomBigInt) IsEmptyOrNil() bool {
	return cbi.String() == "0" || cbi.String() == "<nil>"
}
func HexToBigInt(hexStr string) (*big.Int, error) {
	if hexStr == "" {
		return nil, errors.New("empty hex string")
	}
	bi := new(big.Int)
	_, ok := bi.SetString(hexStr[2:], 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return bi, nil
}

func BigIntToHex(number *big.Int) (string, error) {
	hexStr := number.Text(16)
	if hexStr == "" || hexStr == "<nil>" {
		return "", fmt.Errorf("invalid input")
	}

	return "0x" + hexStr, nil
}
