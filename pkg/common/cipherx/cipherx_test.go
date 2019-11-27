package cipherx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"etop.vn/common/jsonx"
)

func TestCipherX(t *testing.T) {
	cipherx, err := NewCipherx("secretkey")
	assert.NoError(t, err)

	for _, tt := range []struct {
		text string
	}{
		{"dữ liệu mẫu"},
		{"2"},
	} {
		t.Run("Case encrypt string", func(t *testing.T) {
			encryptText, err := cipherx.Encrypt([]byte(tt.text))
			assert.NoError(t, err)

			plainText, err := cipherx.Decrypt(encryptText)
			assert.NoError(t, err)
			if tt.text != string(plainText) {
				t.Errorf("CryptoService.Decrypt() plain = %v, want = %v", plainText, tt.text)
				return
			}
		})
	}

	type testType struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	}

	for _, tt := range []struct {
		test testType
	}{
		{testType{"key1", "value1"}},
		{testType{"key2", "value2"}},
	} {
		t.Run("Case encrypt struct", func(t *testing.T) {
			data, err := json.Marshal(tt.test)
			assert.NoError(t, err)
			encryptText, err := cipherx.Encrypt(data)
			text := string(encryptText)
			assert.NoError(t, err)
			data, err = jsonx.Marshal(text)
			if err != nil {
				panic(err)
			}

			plainText, err := cipherx.Decrypt([]byte(text))
			assert.NoError(t, err)

			var res testType
			err = json.Unmarshal(plainText, &res)
			assert.NoError(t, err)
			if tt.test.Key != res.Key {
				t.Errorf("CryptoService.Decrypt() plain = %v, want = %v", res.Key, tt.test.Key)
				return
			}
		})
	}
}
