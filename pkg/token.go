package pkg

import (
	"crypto/sha256"
	"encoding/hex"
)

func FormatToken(token string) string {
	hash := sha256.New()

	hash.Write([]byte(token))

	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
