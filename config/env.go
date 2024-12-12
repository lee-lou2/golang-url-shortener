package config

import (
	"github.com/joho/godotenv"
	"os"
	"sync"
)

var envOnce sync.Once

// GetEnv 환경 변수 조회
func GetEnv(key string, defaults ...string) string {
	envOnce.Do(func() {
		_ = godotenv.Load()
	})
	if len(defaults) > 0 {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaults[0]
	}
	return os.Getenv(key)
}
