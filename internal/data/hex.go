package data

import (
	"encoding/base64"
	"strings"
)

// ValidHex verifies if str is valid hex
func ValidHex(length int, str string) bool {
	if len(str) != length {
		return false
	}
	for _, c := range str {
		ord := int(c)
		if !(int('0') <= ord && ord <= int('9')) &&
			!(int('a') <= ord && ord <= int('f')) {
			return false
		}
	}
	return true
}

func mbase64(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(data)
	return strings.ReplaceAll(strings.ReplaceAll(s, "/", "_"), "=", ""), nil
}

// ToBase64 convert byte slice to base64 encoded byte slice
func ToBase64(val string) ([]byte, error) {
	s := strings.ReplaceAll(strings.ReplaceAll(val, "/", ""), "_", "/")
	return base64.StdEncoding.DecodeString(s)
}
