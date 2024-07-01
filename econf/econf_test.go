package econf

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Error: %v", err)
	}

	t.Run("LoadConfig", func(t *testing.T) {
		cfg, err := LoadConfig("econf-test")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if cfg.SignatureKeys.PublicKeyPath == "" {
			t.Fatalf("Unexpected keys config: %+v", cfg.SignatureKeys)
		}
		if cfg.SignatureKeys.Version != "econf-test" {
			t.Fatalf("Unexpected version: %+v", cfg.SignatureKeys.Version)
		}
	})
}
