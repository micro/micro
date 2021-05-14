// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/network/tunnel/mucp/crypto_test.go

package mucp

import (
	"bytes"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := []byte("tokenpassphrase")
	gcm, err := newCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("supersecret")

	cipherText, err := Encrypt(gcm, data)
	if err != nil {
		t.Errorf("failed to encrypt data: %v", err)
	}

	// verify the cipherText is not the same as data
	if bytes.Equal(data, cipherText) {
		t.Error("encrypted data are the same as plaintext")
	}
}

func TestDecrypt(t *testing.T) {
	key := []byte("tokenpassphrase")
	gcm, err := newCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("supersecret")

	cipherText, err := Encrypt(gcm, data)
	if err != nil {
		t.Errorf("failed to encrypt data: %v", err)
	}

	plainText, err := Decrypt(gcm, cipherText)
	if err != nil {
		t.Errorf("failed to decrypt data: %v", err)
	}

	// verify the plainText is the same as data
	if !bytes.Equal(data, plainText) {
		t.Error("decrypted data not the same as plaintext")
	}
}
