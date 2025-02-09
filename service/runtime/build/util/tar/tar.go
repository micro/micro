// Package tar is a wrapper around archive/tar, it's used for archiving and unarchiving source from
// and into folders on the host. Because runtime has many cases where it streams source via RPC,
// this package encapsulates the shared logic.
package tar

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Unarchive decodes the source in a tar and writes it to a directory
func Unarchive(src io.Reader, dir string) error {
	tr := tar.NewReader(src)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(dir, hdr.Name)
		bytes, err := ioutil.ReadAll(tr)
		if err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if _, verr := os.Stat(path); os.IsNotExist(verr) {
				err = os.Mkdir(path, os.ModePerm)
			}
		case tar.TypeReg:
			err = ioutil.WriteFile(path, bytes, os.ModePerm)
		default:
			err = fmt.Errorf("Unknown tar header type flag: %v", string(hdr.Typeflag))
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// Archive a local directory into a tar gzip
func Archive(dir string) (io.Reader, error) {
	// Create a tar writer and a buffer to store the archive
	tf := bytes.NewBuffer(nil)
	tw := tar.NewWriter(tf)
	defer tw.Close()

	// walkFn archives each file in the directory
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// get the relative path, e.g. cmd/main.go
		relpath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// skip git files
		if filepath.HasPrefix(relpath, ".git") {
			return nil
		}

		// skip irrelevent files
		if !info.IsDir() && !shouldArchive(relpath) {
			return nil
		}

		// generate and write tar header
		header, err := tar.FileInfoHeader(info, relpath)
		if err != nil {
			return err
		}

		// Since os.FileInfo's Name method only returns the base name of the file it describes, it is
		// necessary to modify Header.Name to provide the full path name of the file. See:
		// https://golang.org/src/archive/tar/common.go?s=22088:22153#L626
		header.Name = relpath

		// write the header to the archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// there is no body if it's a directory
		if info.IsDir() {
			return nil
		}

		// read the contents of the file
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// write the contents of the file to the tar
		_, err = tw.Write([]byte(bytes))
		return err
	}

	// Add the files to the archive
	if err := filepath.Walk(dir, walkFn); err != nil {
		return nil, err
	}

	// Return the archive
	return tf, nil
}

// shouldArchive is a helper func which indicates if a file should be archived. TODO: implement a
// smarter check which just excludes executables
func shouldArchive(file string) bool {
	if filepath.HasPrefix(file, "vendor") {
		return true
	}
	if strings.HasSuffix(file, ".go") {
		return true
	}
	if strings.HasSuffix(file, "go.sum") {
		return true
	}
	if strings.HasSuffix(file, "go.mod") {
		return true
	}
	if strings.HasSuffix(file, ".txt") {
		return true
	}
	return false
}
