package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v3/runtime/local/source/git"
	"github.com/micro/go-micro/v3/util/tar"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/store"
)

const (
	sourcePrefix = "source://"
)

// WriteSource to the blob store. Returns the key and an error if one occurs.
func WriteSource(namespace string, src io.Reader) (string, error) {
	// generate a key to use as the source
	key := sourcePrefix + uuid.New().String()

	// write it to the blob store
	err := store.DefaultBlobStore.Write(key, src)
	if err != nil {
		return "", fmt.Errorf("Error writing source to blob store: %v", err)
	}

	// return the key
	return key, nil
}

// ReadSource takes the source of a service passed in runtime.Create/Update, this value could be
// the upload ID e.g. source://foo-bar or a git remote. ReadSource will load this source from the
// respective location. The resuling source will archived in a tar gzip.
func ReadSource(srv *runtime.Service, secrets map[string]string, namespace string) (io.Reader, error) {
	// validate the arguments
	if len(srv.Source) == 0 {
		return nil, fmt.Errorf("Missing source")
	}
	if len(namespace) == 0 {
		return nil, fmt.Errorf("Missing namespace")
	}

	// source was previously uploaded to the blob store
	if strings.HasPrefix(srv.Source, sourcePrefix) {
		return store.DefaultBlobStore.Read(srv.Source)
	}

	// by process of elimination, the source is a git remote. we'll now load it. the go-micro git package
	// loads the git source a directory so we'll create a temp directory to use
	tmpDir, err := ioutil.TempDir(os.TempDir(), "source")
	if err != nil {
		return nil, err
	}

	// CheckoutSource requires secrets incase they contain any Git credentials (todo: replace this
	// with an option).
	src, err := git.ParseSource(srv.Source)
	if err != nil {
		return nil, err
	}
	if err := git.CheckoutSource(tmpDir, src, secrets); err != nil {
		return nil, err
	}

	// Archive the source and then return
	return tar.Archive(tmpDir)
}
