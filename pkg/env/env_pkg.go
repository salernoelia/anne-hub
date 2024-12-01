package env

import (
	"log"
	"os"
)

func GetEnvOrFatal(key string, defaultValue ...string) string {
    value := os.Getenv(key)
    if value == "" {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        log.Fatalf("Environment variable %s is not set.", key)
    }
    return value
}