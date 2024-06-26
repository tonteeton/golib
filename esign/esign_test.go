package esign

import (
	"bytes"
	"crypto/ed25519"
	"github.com/tonteeton/golib/econf"
	"os"
	"testing"
)

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

func TestSignatureKey(t *testing.T) {
	config := econf.KeysConfig{
		PublicKeyPath:  "key.pub",
		PrivateKeyPath: "key.priv.enc",
		SealedDatePath: "created.enc",
		Version:        "test",
	}

	t.Run("Keys loaded correctly", func(t *testing.T) {
		defer setupTest(t, config)()
		_, err := GetSignatureKey(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	t.Run("Key loaded and reused", func(t *testing.T) {
		defer setupTest(t, config)()

		key1, err := GetSignatureKey(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if _, err = os.Stat(config.PrivateKeyPath); err != nil {
			t.Fatalf("Key was not saved")
		}

		key2, err := GetSignatureKey(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if !key1.PrivateKey.Equal(key2.PrivateKey) {
			t.Fatalf("Key was not reused")
		}

		err = os.Remove(config.PrivateKeyPath)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		key3, err := GetSignatureKey(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if key1.PrivateKey.Equal(key3.PrivateKey) {
			t.Fatalf("Private keys always same, %x %x", key1.PrivateKey, key3.PrivateKey)
		}
	})

	t.Run("Public key returned", func(t *testing.T) {
		defer setupTest(t, config)()
		key, err := GetSignatureKey(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		publicKey := key.GetPublicKey()
		if len(publicKey) != ed25519.PublicKeySize {
			t.Fatalf("Key of unexpected length")
		}
	})

	t.Run("Key saved", func(t *testing.T) {
		defer setupTest(t, config)()

		var priv [32]byte
		for i := range priv {
			priv[i] = 0xaa
		}

		var pub [32]byte
		for i := range pub {
			pub[i] = 0xbb
		}

		err := SaveSignatureKey(config, append(priv[:], pub[:]...))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		key, err := GetSignatureKey(config)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		_ = key

		if !bytes.Equal(key.GetPublicKey(), pub[:]) {
			t.Fatalf("Wrong public key loaded: %#v", key.GetPublicKey())
		}

		if !bytes.Equal(key.GetPrivateKey(), append(priv[:], pub[:]...)) {
			t.Fatalf("Wrong private key loaded: %#v", key.GetPrivateKey())
		}

	})
}
