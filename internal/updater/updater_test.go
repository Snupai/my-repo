package updater

import (
	"testing"
)

func TestGetBinaryName(t *testing.T) {
	binaryName := getBinaryName()
	if binaryName == "" {
		t.Error("Binary name should not be empty")
	}

	// Should contain "cngt-cli"
	if !contains(binaryName, "cngt-cli") {
		t.Error("Binary name should contain 'cngt-cli'")
	}
}

func TestCheckForUpdates(t *testing.T) {
	// This test requires internet connection, so we'll make it optional
	t.Skip("Skipping CheckForUpdates test - requires internet connection")
	
	release, hasUpdate, err := CheckForUpdates()
	if err != nil {
		t.Logf("Error checking for updates: %v", err)
		return
	}

	if release == nil {
		t.Error("Release should not be nil")
	}

	t.Logf("Has update: %v", hasUpdate)
	if release != nil {
		t.Logf("Latest version: %s", release.TagName)
	}
}

func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}