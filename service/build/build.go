package build

import (
	"io"
)

// DefaultBuilder implementation. Note: we don't set the client here as that would result in a
// circular dependancy. This isn't an issue with other interfaces as they're normally defined in
// go-micro. Profiles should configure this builder but clients of this package should handle the
// nil value case.
var DefaultBuilder Builder

// Builder is an interface for building packages
type Builder interface {
	// Build a package
	Build(src io.Reader, opts ...Option) (io.Reader, error)
}
