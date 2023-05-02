package utils

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestHexToBigInt(t *testing.T) {
	type args struct {
		hexStr string
	}
	tests := []struct {
		name     string
		args     args
		expected *big.Int
		hasError bool
	}{
		{
			name: "test_hex_to_big_int_success",
			args: args{
				hexStr: "0xbc614e",
			},
			expected: big.NewInt(12345678),
			hasError: false,
		},
		{
			name: "test_hex_to_big_int_fail_when_hex_string_is_empty",
			args: args{
				hexStr: "",
			},
			expected: nil,
			hasError: true,
		},
		{
			name: "test_hex_to_big_int_fail_invalid_hex_string",
			args: args{
				hexStr: "zzzzz",
			},
			expected: nil,
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := HexToBigInt(tt.args.hexStr)

			if tt.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBigIntToHex(t *testing.T) {
	type args struct {
		number *big.Int
	}
	tests := []struct {
		name     string
		args     args
		expected string
		hasError bool
	}{
		{
			name: "test_big_int_to_hex_success",
			args: args{
				number: big.NewInt(12345678),
			},
			expected: "0xbc614e",
			hasError: false,
		},
		{
			name: "test_big_int_to_hex_fail_big_int_is_nil",
			args: args{
				number: nil,
			},
			expected: "",
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := BigIntToHex(tt.args.number)

			if tt.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_IsEmptyOrNil_whenEmpty(t *testing.T) {
	empty := new(big.Int)
	customEmpty := CustomBigInt{empty}

	assert.True(t, customEmpty.IsEmptyOrNil())
}

func Test_IsEmptyOrNil_whenNil(t *testing.T) {
	var bigIntNil *big.Int
	customNil := CustomBigInt{bigIntNil}

	assert.True(t, customNil.IsEmptyOrNil())
}
