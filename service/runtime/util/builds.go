package util

import (
	"io"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/store"
)

const (
	buildPrefix = "build://"
)

// WriteBuild to the blob store. Returns the key and an error if one occurs.
func WriteBuild(namespace string, src io.Reader) (string, error) {
	// generate a key to identify the build
	key := buildPrefix + uuid.New().String()

	// write it to the blob store
	err := store.DefaultBlobStore.Write(key, src)
	if err != nil {
		return "", err
	}

	// return the key
	return key, nil
}
