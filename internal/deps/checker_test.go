package deps

import (
	"testing"
)

func TestRequiredPackages(t *testing.T) {
	if len(requiredPackages) == 0 {
		t.Error("Required packages should not be empty")
	}

	expectedPackages := []string{"termcolor", "mido", "colorama", "cryptography"}
	for _, expected := range expectedPackages {
		found := false
		for _, pkg := range requiredPackages {
			if pkg == expected || (len(pkg) > len(expected) && pkg[:len(expected)] == expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected package %s not found in required packages", expected)
		}
	}
}

func TestIsPythonInstalled(t *testing.T) {
	// This test depends on the environment, so we'll just check that it doesn't panic
	result := isPythonInstalled()
	t.Logf("Python installed: %v", result)
}

func TestIsUvAvailable(t *testing.T) {
	// This test depends on the environment, so we'll just check that it doesn't panic
	result := isUvAvailable()
	t.Logf("uv available: %v", result)
}