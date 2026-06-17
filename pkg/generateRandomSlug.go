package pkg

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomSlug(length int) (string, error) {
	slug := make([]byte, length)

	for i := range slug {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		slug[i] = charset[n.Int64()]
	}

	return string(slug), nil
}
