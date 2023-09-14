package common

import (
	"fmt"
	"math/big"
	"strings"
)

// TODO 인터페이스로 멀티체인 추가

func HexToDecimal(str string) (string, error) {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}

	num := big.NewInt(0)
	num, success := num.SetString(str, 16)
	if !success {
		return "", fmt.Errorf("failed to convert hex string to big.Int")
	}

	return num.String(), nil
}
