// Package zip is a wrapper around archive/zip, it's used for archiving and unarchiving source from
// and into folders on the host. Because runtime has many cases where it streams source via RPC,
// this package encapsulates the shared logic.
package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Unarchive decodes the source in a zip and writes it to a directory
func Unarchive(src io.Reader, dir string) error {
	// create a new buffer with the source, this is required because zip.NewReader takes a io.ReaderAt
	// and not an io.Reader
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, src)
	if err != nil {
		return err
	}

	// create the zip
	reader := bytes.NewReader(buff.Bytes())
	zip, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	// write the files in the zip to our tmp dir
	for _, f := range zip.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}

		path := filepath.Join(dir, f.Name)
		if err := ioutil.WriteFile(path, bytes, os.ModePerm); err != nil {
			return err
		}

		if err := rc.Close(); err != nil {
			return err
		}
	}

	return nil
}
