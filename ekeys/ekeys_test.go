package ekeys

import (
	"bytes"
	"crypto/rand"
	"github.com/edgelesssys/ego/ecrypto"
	"github.com/tonteeton/golib/econf"
	"os"
	"strings"
	"testing"
)

func generateRandomKey() (publicKey []byte, privateKey []byte, err error) {
	var expectedPub [32]byte
	var expectedPriv [32]byte
	if _, err := rand.Read(expectedPub[:]); err != nil {
		return nil, nil, err
	}
	if _, err := rand.Read(expectedPriv[:]); err != nil {
		return nil, nil, err
	}
	return expectedPub[:], expectedPriv[:], nil
}

func setupTest(t *testing.T, config econf.KeysConfig) func() {
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Error: %v", err)
	}
	if _, err := os.Stat(config.PrivateKeyPath); err == nil {
		err = os.Remove(config.PrivateKeyPath)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	}

	return func() {

	}
}

func TestLoadConfig(t *testing.T) {
	config := econf.KeysConfig{
		PublicKeyPath:  "key.pub",
		PrivateKeyPath: "key.priv.enc",
		SealedDatePath: "created.enc",
		Version:        "test",
	}
	keys := KeyManager{config}

	t.Run("Write and read encrypted file", func(t *testing.T) {
		setupTest(t, config)
		err := WriteEncryptedFile("test.enc", []byte("testdata"), []byte("test"))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		var data []byte
		data, err = os.ReadFile("test.enc")
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if strings.Contains(string(data), "testdata") {
			t.Fatalf("Data was not encrypted: '%s'", data)
		}

		var decrypted []byte
		decrypted, err = ReadEncryptedFile("test.enc", []byte("test"))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if string(decrypted) != "testdata" {
			t.Fatalf("Unexpected decrypted data: '%s'", decrypted)
		}

	})

	t.Run("Write and read file encrypted with a Product key", func(t *testing.T) {
		setupTest(t, config)
		err := WriteEncryptedFile("test.enc", []byte("testdata"), []byte("test"), ecrypto.SealWithProductKey)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		var data []byte
		data, err = os.ReadFile("test.enc")
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if strings.Contains(string(data), "testdata") {
			t.Fatalf("Data was not encrypted: '%s'", data)
		}

		var decrypted []byte
		decrypted, err = ReadEncryptedFile("test.enc", []byte("test"))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if string(decrypted) != "testdata" {
			t.Fatalf("Unexpected decrypted data: '%s'", decrypted)
		}
	})

	t.Run("Create and load keys", func(t *testing.T) {
		setupTest(t, config)
		_, err := keys.load()
		if err == nil || string(err.Error()) != "failed to read creation info" {
			t.Fatalf("Error on reading creation info expeted, got: %v", err)
		}

		err = keys.Save([]byte("testpub"), []byte("testpriv"))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		key, err := keys.load()
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if !bytes.Equal(key, []byte("testpriv")) {
			t.Fatalf("Unexpected key loaded: %#v", key)
		}
	})

	t.Run("Key loaded and reused", func(t *testing.T) {
		setupTest(t, config)
		key1, err := keys.GetPrivateKey(generateRandomKey)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if _, err = os.Stat(config.PrivateKeyPath); err != nil {
			t.Fatalf("Key was not saved")
		}

		key2, err := keys.GetPrivateKey(generateRandomKey)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if !bytes.Equal(key1, key2) {
			t.Fatalf("Key was not reused")
		}

		err = os.Remove(config.PrivateKeyPath)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		key3, err := keys.GetPrivateKey(generateRandomKey)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if bytes.Equal(key1, key3) {
			t.Fatalf("Private keys always same, %x %x", key1, key3)
		}
	})

	t.Run("Invalid Key not used", func(t *testing.T) {
		setupTest(t, config)
		_, err := keys.GetPrivateKey(generateRandomKey)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if err := os.WriteFile(config.PrivateKeyPath, []byte("modified"), 0644); err != nil {
			t.Fatalf("Error writing file: %v", err)
		}

		_, err = keys.GetPrivateKey(generateRandomKey)
		if err == nil || !strings.Contains(err.Error(), "private key") {
			t.Fatalf("Expected 'unexpected keys creation date' error, got: %v", err)
		}
	})
}
