package utils

import (
	"fmt"
)

func ConvertToEmployeeCode(n int) (string, error) {
	if n >= 26*26*10 {
		return "", fmt.Errorf("รหัสเกินขีดจำกัด (%d)", 26*26*10)
	}

	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits := []rune("0123456789")

	// คำนวณลำดับ
	letterIndex := n / 10
	digitIndex := n % 10

	first := letters[letterIndex/26]
	second := letters[letterIndex%26]
	digit := digits[digitIndex]

	return fmt.Sprintf("#%c%c%c", first, second, digit), nil
}
