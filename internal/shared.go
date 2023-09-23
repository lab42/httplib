package internal

import (
	"net/http"
	"strings"
)

// obfuscateSensitiveData obfuscates sensitive data based on configured sensitive keys.
func ObfuscateSensitiveData(key string, values []string, sensitiveKeys []string) []string {
	var obfuscatedValues []string
	if IsSensitiveField(key, sensitiveKeys) {
		// Obfuscate the sensitive data with asterisks
		obfuscatedValues = append(obfuscatedValues, "********")
	} else {
		obfuscatedValues = append(obfuscatedValues, values...)
	}
	return obfuscatedValues
}

// isSensitiveField checks if a field (identified by key) contains sensitive data based on sensitive keys.
func IsSensitiveField(key string, sensitiveKeys []string) bool {
	// Check if the key matches any of the sensitive keys
	for _, sensitiveKey := range sensitiveKeys {
		if strings.EqualFold(key, sensitiveKey) {
			return true
		}
	}
	return false
}

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	// You can write a response or perform any other desired actions here.
	// In most cases, it's left empty or just responds with a status code.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
