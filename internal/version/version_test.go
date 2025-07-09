package version

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("Version should not be empty")
	}
}

func TestGetFullVersion(t *testing.T) {
	fullVersion := GetFullVersion()
	if fullVersion == "" {
		t.Error("Full version should not be empty")
	}
	
	// Should contain version, commit, and build time
	if !contains(fullVersion, "commit:") {
		t.Error("Full version should contain commit info")
	}
	
	if !contains(fullVersion, "built:") {
		t.Error("Full version should contain build time")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		   len(s) > len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}