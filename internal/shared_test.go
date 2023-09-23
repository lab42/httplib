package internal_test

import (
	"testing"

	"github.com/lab42/httplib/internal"
)

func TestObfuscateSensitiveData(t *testing.T) {
	// Test case 1: Key is sensitive, values should be obfuscated.
	sensitiveKeys := []string{"password", "ssn"}
	key := "password"
	values := []string{"mysecret"}

	obfuscatedValues := internal.ObfuscateSensitiveData(key, values, sensitiveKeys)
	expectedObfuscatedValues := []string{"********"}
	if !stringSlicesEqual(obfuscatedValues, expectedObfuscatedValues) {
		t.Errorf("Expected obfuscated values %v, but got %v", expectedObfuscatedValues, obfuscatedValues)
	}

	// Test case 2: Key is not sensitive, values should remain unchanged.
	key = "username"
	values = []string{"john.doe"}

	obfuscatedValues = internal.ObfuscateSensitiveData(key, values, sensitiveKeys)
	if !stringSlicesEqual(obfuscatedValues, values) {
		t.Errorf("Expected values %v, but got %v", values, obfuscatedValues)
	}
}

func TestIsSensitiveField(t *testing.T) {
	// Test case 1: Key matches a sensitive key.
	sensitiveKeys := []string{"password", "ssn"}
	key := "password"

	result := internal.IsSensitiveField(key, sensitiveKeys)
	if !result {
		t.Errorf("Expected IsSensitiveField to return true for key '%s', but got false", key)
	}

	// Test case 2: Key does not match any sensitive key.
	key = "username"

	result = internal.IsSensitiveField(key, sensitiveKeys)
	if result {
		t.Errorf("Expected IsSensitiveField to return false for key '%s', but got true", key)
	}
}

// Helper function to check if two string slices are equal.
func stringSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
