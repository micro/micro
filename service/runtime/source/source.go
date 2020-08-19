// package source
package source

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/internal/git"
)

var (
	// The source dir
	Dir = filepath.Join(os.TempDir(), "micro", "source")
)

// Checkout source code if needed
func Checkout(s *runtime.Service) error {
	// Runtime service like config have no source.
	// Skip checkout in that case
	if len(s.Source) == 0 {
		return nil
	}
	// @todo make this come from config
	cpath := filepath.Join(Dir, s.Source)

	path := strings.ReplaceAll(cpath, ".tar.gz", "")

	// unpack the tarball if it exists
	if ex, _ := exists(cpath); ex {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
		err = git.Uncompress(cpath, path)
		if err != nil {
			return err
		}
		s.Source = path
		return nil
	}

	// if the tarball does not exit try checkout the source code
	source, err := git.ParseSourceLocal("", s.Source)
	if err != nil {
		return err
	}
	source.Ref = s.Version

	err = git.CheckoutSource(os.TempDir(), source)
	if err != nil {
		return err
	}
	s.Source = source.FullPath
	return nil
}
