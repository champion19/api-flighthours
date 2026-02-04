package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindModuleRoot(t *testing.T) {
	t.Run("finds module root from current directory", func(t *testing.T) {
		root, err := FindModuleRoot()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if root == "" {
			t.Error("expected non-empty root path")
		}

		// Verify go.mod exists in the returned path
		goModPath := filepath.Join(root, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			t.Errorf("go.mod not found at %s", goModPath)
		}
	})
}
