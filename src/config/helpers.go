//helpers.go

package config

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getString(jsonVal, envKey, defaultVal string) string {
	if jsonVal != "" {
		return jsonVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		return envVal
	}
	return defaultVal
}

func getStringSlice(jsonValue []string, envKey string, fallback []string) []string {
	if len(jsonValue) > 0 {
		return jsonValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		// Split the environment variable by commas, trim whitespace
		parts := strings.Split(envValue, ",")
		var result []string
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return fallback
}

func getInt(jsonVal int, envKey string, defaultVal int) int {
	if jsonVal != 0 {
		return jsonVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		if val, err := strconv.Atoi(envVal); err == nil {
			return val
		}
	}
	return defaultVal
}
func getInt64(jsonVal int64, envKey string, defaultVal int64) int64 {
	if jsonVal != 0 {
		return jsonVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		if val, err := strconv.ParseInt(envVal, 10, 64); err == nil {
			return val
		}
	}
	return defaultVal
}

func getBool(jsonVal *bool, envKey string, defaultVal bool) bool {
	if jsonVal != nil {
		return *jsonVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			return val
		}
	}
	return defaultVal
}

// getUsers retrieves a map[string]string with JSON -> env -> default hierarchy
func getUsers(jsonValue map[string]string, envKey string, defaultValue map[string]string) map[string]string {
	if jsonValue != nil {
		return jsonValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		// Expect env var as "user1:hash1,user2:hash2"
		users := make(map[string]string)
		pairs := strings.Split(envValue, ",")
		for _, pair := range pairs {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				users[parts[0]] = parts[1]
			}
		}
		if len(users) > 0 {
			return users
		}
	}
	return defaultValue
}

func getDuration(jsonVal time.Duration, envKey string, defaultVal time.Duration) time.Duration {
	if jsonVal != 0 {
		return jsonVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		if val, err := time.ParseDuration(envVal); err == nil {
			return val
		}
	}
	return defaultVal
}

func generateJwtKey() string {

	// ensure we return JwtKey if it's set
	if len(JwtKey) > 0 {
		return JwtKey
	}
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Failed to generate JWT key, using fallback")
		return "i-am-a-fallback-32-byte-secret-key!!"
	}
	return base64.RawURLEncoding.EncodeToString(key)
}

// runDeferredAction runs the provided action after unlocking the mutex
func runDeferredAction(action DeferredAction) {
	if action != nil {
		go action() // Run in a goroutine to ensure non-blocking execution
	}
}

// Getter and Setter for the bundled assets. Could be in getters.go and setters.go, but made more sense to keep these special cases out of there.
func GetV1UIFS() embed.FS {
	return V1UIFS
}

func GetV2UIFS() embed.FS {
	return V2UIFS
}

func GetTWOBOXFS() embed.FS {
	return TWOBOXFS
}

func SetTWOBOXFS(twoboxFS embed.FS) {
	TWOBOXFS = twoboxFS
}

func SetV1UIFS(v1uiFS embed.FS) {
	V1UIFS = v1uiFS
}

func SetV2UIFS(v2uiFS embed.FS) {
	V2UIFS = v2uiFS
}
