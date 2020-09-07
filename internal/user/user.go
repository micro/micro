package user

import (
	"log"
	"os/user"
	"path/filepath"
	"runtime/debug"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/micro/micro/v3/service/logger"
)

var (
	Dir  = ""
	path = ".micro"
)

func init() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	Dir = filepath.Join(user.HomeDir, path)
}

func GetKeys() ([]byte, []byte, error) {
	debug.PrintStack()
	privKey := filepath.Join(Dir, "id_rsa")
	pubKey := filepath.Join(Dir, "id_rsa.pub")

	if !fileExists(privKey) || !fileExists(pubKey) {
		err := setupKeys(privKey, pubKey)
		if err != nil {
			logger.Fatalf("Error setting up keys: %v", err)
		}
	}
	privDat, err := ioutil.ReadFile(privKey)
	if err != nil {
		return nil, nil, err
	}
	pubDat, err := ioutil.ReadFile(pubKey)
	if err != nil {
		return nil, nil, err
	}
	return privDat, pubDat, nil
}

func setupKeys(privKey, pubKey string) error {
	bitSize := 4096
	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return err
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	privateKeyBytes, err := encodePrivateKeyToPEM(privateKey)
	if err != nil {
		return err
	}

	err = writeKeyToFile(privateKeyBytes, privKey)
	if err != nil {
		return err
	}

	err = writeKeyToFile([]byte(publicKeyBytes), pubKey)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// taken from https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM, nil
}

func generatePublicKey(publickey *rsa.PublicKey) ([]byte, error) {
	pubDER, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return nil, err
	}
	pubBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDER,
	}
	pubPEM := pem.EncodeToMemory(&pubBlock)

	return pubPEM, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	file, err := os.OpenFile(saveFileTo, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Close()

	err = ioutil.WriteFile(saveFileTo, []byte(base64.StdEncoding.EncodeToString(keyBytes)), 0600)
	if err != nil {
		return err
	}

	return nil
}
