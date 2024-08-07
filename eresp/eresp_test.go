package eresp

import (
	"encoding/base64"
	"encoding/json"
	"github.com/tonteeton/golib/econf"
	"github.com/tonteeton/golib/esign"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"os"
	"testing"
)

func setupTest(t *testing.T) func() {
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Error: %v", err)
	}
	return func() {}
}

func TestEnclaveResponse(t *testing.T) {
	secretKey, err := base64.StdEncoding.DecodeString(
		"yMJNiUZf3kMeEkQ+0r57+Ou8DEfOKmNC/BCN9c2TfPc5PICixeaQ8vlV/79OARLthRMyTOXEVDU16/1JY3BP1Q==",
	)
	if err != nil {
		t.Fatalf("Error decoding secret key: %v", err)
	}
	boc64 := "te6cckEBAQEAMgAAYAAAAABmOjrBAAAAAHJxYCMAAAAAAAABWQAAABMVr91EAAAAAAAABh4AAAAAAAAq11siUa4="
	data, err := base64.StdEncoding.DecodeString(boc64)
	if err != nil {
		t.Fatalf("Error decoding base64: %v", err)
	}
	payload, err := cell.FromBOC(data)
	if err != nil {
		t.Fatalf("Error building cell from BOC: %v", err)
	}

	t.Run("ValidInputs", func(t *testing.T) {
		setupTest(t)
		got, err := newEnclaveResponse(payload, secretKey)
		if err != nil {
			t.Errorf("Error building Enclave response: %v", err)
		}

		expectedPayload := boc64
		if got.Payload != expectedPayload {
			t.Errorf("Unexpected payload. Got: %s, Expected: %s", got.Payload, expectedPayload)
		}

		expectedHash := "KWraQp7R+lYAaGw9VqJnMeKcar9q+mKtudCST/4h3GY="
		if got.Hash != expectedHash {
			t.Errorf("Unexpected hash. Got: %s, Expected: %s", got.Hash, expectedHash)
		}

		expectedSignature := "Id3NO8Tbq4ZFcZ1mp4gr78g7+SgmHuCdTSSBXmzXYy7u3W/UPisnTsE7CuDUATiaOFnE208w1fyb8+s6BM/0BA=="
		if got.Signature != expectedSignature {
			t.Errorf("Unexpected signature. Got: %s, Expected: %s", got.Signature, expectedSignature)
		}
	})

	t.Run("InvalidSecretKey", func(t *testing.T) {
		setupTest(t)
		_, err := newEnclaveResponse(payload, []byte("invalidsecretkey"))
		if err != nil {
			t.Error("Error expected for invalid secret key")
		}
	})

}

func TestEnclaveResponseSaveMethod(t *testing.T) {

	t.Run("ResponseSavedToJson", func(t *testing.T) {
		setupTest(t)
		response := EnclaveResponse{
			Signature: "signature",
			Payload:   "payload",
			Hash:      "hash",
		}

		err := response.save("response.json")
		if err != nil {
			t.Fatalf("Error saving EnclaveResponse to JSON: %v", err)
		}

		data, _ := os.ReadFile("response.json")
		var savedResponse EnclaveResponse
		err = json.Unmarshal(data, &savedResponse)
		if err != nil {
			t.Fatalf("Error unmarshaling saved JSON data: %v", err)
		}
	})
}

func TestSaveResponse(t *testing.T) {

	t.Run("SaveResponseNoError", func(t *testing.T) {
		setupTest(t)
		responseCfg := Config{
			Response: econf.ResponseConfig{
				ResponsePath: "response1.json",
			},
			SignatureKeys: econf.KeysConfig{
				PublicKeyPath:  "key.pub",
				PrivateKeyPath: "key.priv.enc",
				SealedDatePath: "created.enc",
				Version:        "test",
			},
		}
		err := SaveResponse(responseCfg, cell.BeginCell().EndCell())
		if err != nil {
			t.Fatalf("Error on SaveResponse: %v", err)
		}
		_, err = os.Stat("response1.json")
		if err != nil {
			t.Fatalf("Response file was not created: %v.", err)
		}

	})
}

func TestPackResponseToCell(t *testing.T) {
	t.Run("Response packed to cell as expected", func(t *testing.T) {
		setupTest(t)

		responseCfg := Config{
			Response: econf.ResponseConfig{
				ResponsePath: "response1.json",
			},
			SignatureKeys: econf.KeysConfig{
				PublicKeyPath:  "key.pub",
				PrivateKeyPath: "key.priv.enc",
				SealedDatePath: "created.enc",
				Version:        "test",
			},
		}

		secretKey, err := base64.StdEncoding.DecodeString(
			"yMJNiUZf3kMeEkQ+0r57+Ou8DEfOKmNC/BCN9c2TfPc5PICixeaQ8vlV/79OARLthRMyTOXEVDU16/1JY3BP1Q==",
		)
		if err != nil {
			t.Fatalf("Error decoding signature key: %v", err)
		}
		err = esign.SaveSignatureKey(
			responseCfg.SignatureKeys,
			secretKey,
		)
		if err != nil {
			t.Fatalf("Error saving signature: %v", err)
		}

		boc64 := "te6cckEBAQEAMgAAYAAAAABmalYRAAAAAHJxYCMAAAAAAAAC9AAAAAnwlP1UAAAAAAAAA2QAAAAAAAArd2S53VY="
		data, err := base64.StdEncoding.DecodeString(boc64)
		if err != nil {
			t.Fatalf("Error decoding base64: %v", err)
		}

		payloadCell, err := cell.FromBOC(data)
		if err != nil {
			t.Fatalf("Error building cell from BOC: %v", err)
		}

		resultCell, err := PackResponseToCell(responseCfg, payloadCell, 0x9f89304e)
		if err != nil {
			t.Errorf("Error packing Enclave response to cell: %v", err)
		}

		expected := "te6cckEBAgEAeQABaJ+JME4AAAAAZmpWEQAAAABycWAjAAAAAAAAAvQAAAAJ8JT9VAAAAAAAAANkAAAAAAAAK3cBAIDQxyqnFZq6P51cgXzD37pklWI2NSRjpwoaKWQY3SkrV59otUulbVoU7JMhyM3LYU3u4k/prBCqNkK6G2MPSysIfrV1mg=="

		got := base64.StdEncoding.EncodeToString(resultCell.ToBOC())
		if got != expected {
			t.Errorf("Unexpected cell. Got: %s, Expected: %s", got, expected)
		}
	})

}
