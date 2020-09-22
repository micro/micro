package util

import (
	"io"

	"github.com/google/uuid"
	gostore "github.com/micro/go-micro/v3/store"
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
	err := store.DefaultBlobStore.Write(key, src, gostore.BlobNamespace(blobNamespacePrefix+namespace))
	if err != nil {
		return "", err
	}

	// return the key
	return key, nil
}
