package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GetRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func CheckIfSliceContains[T comparable](slice []T, ele T) bool {

	hit := false

	for i := 0; i < len(slice); i++ {

		if slice[i] == ele {

			hit = true

			return hit
		}

	}

	return hit

}
