package infra

import (
	"fmt"
	"os"
)

// Noop is a task that prints the stage it is on, but otherwise does nothing
type Noop struct {
	ID   string
	Name string
}

// Validate prints Validating
func (n *Noop) Validate() error {
	_, err := fmt.Fprintf(os.Stderr, "[%s] Validating (no-op)\n", n.Name)
	return err
}

// Plan prints Planning
func (n *Noop) Plan() error {
	_, err := fmt.Fprintf(os.Stderr, "[%s] Planning (no-op)\n", n.Name)
	return err
}

// Apply prints Applying
func (n *Noop) Apply() error {
	_, err := fmt.Fprintf(os.Stderr, "[%s] Applying (no-op)\n", n.Name)
	return err
}

// Finalise prints Finalising
func (n *Noop) Finalise() error {
	_, err := fmt.Fprintf(os.Stderr, "[%s] Finalising (no-op)\n", n.Name)
	return err
}

// Destroy prints Destroying
func (n *Noop) Destroy() error {
	_, err := fmt.Fprintf(os.Stderr, "[%s] Destroying (no-op)\n", n.Name)
	return err
}
