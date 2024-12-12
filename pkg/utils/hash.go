package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fit/config"
)

// SHA256 SHA-256 해시 생성
func SHA256(data string) string {
	dataWithSalt := data + config.GetEnv("SHA256_SALT")
	hash := sha256.New()
	hash.Write([]byte(dataWithSalt))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString
}
