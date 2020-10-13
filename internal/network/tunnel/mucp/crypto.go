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
// Original source: github.com/micro/go-micro/v3/network/tunnel/mucp/crypto.go

package mucp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"github.com/micro/micro/v3/internal/network/tunnel"
	"github.com/oxtoacart/bpool"
)

var (
	// the local buffer pool
	// gcmStandardNonceSize from crypto/cipher/gcm.go is 12 bytes
	// 100 - is max size of pool
	noncePool = bpool.NewBytePool(100, 12)
)

// hash hahes the data into 32 bytes key and returns it
// hash uses sha256 underneath to hash the supplied key
func hash(key []byte) []byte {
	sum := sha256.Sum256(key)
	return sum[:]
}

// Encrypt encrypts data and returns the encrypted data
func Encrypt(gcm cipher.AEAD, data []byte) ([]byte, error) {
	var err error

	// get new byte array the size of the nonce from pool
	// NOTE: we might use smaller nonce size in the future
	nonce := noncePool.Get()
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	defer noncePool.Put(nonce)

	// NOTE: we prepend the nonce to the payload
	// we need to do this as we need the same nonce
	// to decrypt the payload when receiving it
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts the payload and returns the decrypted data
func newCipher(key []byte) (cipher.AEAD, error) {
	var err error

	// generate a new AES cipher using our 32 byte key for decrypting the message
	c, err := aes.NewCipher(hash(key))
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

func Decrypt(gcm cipher.AEAD, data []byte) ([]byte, error) {
	var err error

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return nil, tunnel.ErrDecryptingData
	}

	// NOTE: we need to parse out nonce from the payload
	// we prepend the nonce to every encrypted payload
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	ciphertext, err = gcm.Open(ciphertext[:0], nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}
