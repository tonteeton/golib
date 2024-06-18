// Package econf provides configuration management for enclaved applications.
// It defines and initializes various configuration sections, ensuring the necessary
// 'mount' directory exists for storing application data, reports, and outputs.

package econf

import (
	"os"
)

// ResponseConfig represents the configuration for the application response file path.
type ResponseConfig struct {
	ResponsePath string
}

// KeysConfig represents configuration for managing keys.
type KeysConfig struct {
	PublicKeyPath  string
	PrivateKeyPath string
	SealedDatePath string
	Version        string
}

// ReportsConfig represents configuration for signed reports.
type ReportsConfig struct {
	PublicKeysPath       string
	SignatureRequestPath string
	SignatureImportPath  string
	SignatureExportPath  string
}

// Config aggregates all configuration sections.
type Config struct {
	Response       ResponseConfig
	Reports        ReportsConfig
	SignatureKeys  KeysConfig
	EncryptionKeys KeysConfig
}

// LoadConfig initializes a Config struct.
// It also ensures the necessary 'mount' directory exists.
func LoadConfig(appVersion string) (*Config, error) {
	cfg := Config{}
	cfg.Response = ResponseConfig{
		ResponsePath: "mount/response.json",
	}
	cfg.Reports = ReportsConfig{
		PublicKeysPath:       "mount/report_keys.pub",
		SignatureRequestPath: "mount/report_signature_request.pub",
		SignatureImportPath:  "mount/report_signature_import.enc",
		SignatureExportPath:  "mount/report_signature_export.enc",
	}
	cfg.SignatureKeys = KeysConfig{
		PublicKeyPath:  "mount/signature_key.pub",
		PrivateKeyPath: "mount/signature_key.priv.enc",
		SealedDatePath: "mount/signature_created.enc",
		Version:        appVersion,
	}
	cfg.EncryptionKeys = KeysConfig{
		PublicKeyPath:  "mount/box_key.pub",
		PrivateKeyPath: "mount/box_key.priv.enc",
		SealedDatePath: "mount/box_created.enc",
		Version:        appVersion,
	}

	// Ensure the 'mount' directory exists
	if err := os.MkdirAll("mount", 0700); err != nil {
		return nil, err
	}

	return &cfg, nil
}
