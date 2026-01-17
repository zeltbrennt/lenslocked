package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("bytes: didn't read enough bytes")
	}
	return b, nil
}

func generateString(n int) (string, error) {
	b, err := generateBytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func RandomString(n int) (string, error) {
	return generateString(n)
}
