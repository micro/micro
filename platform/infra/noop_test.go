package infra

import (
	"testing"
)

func TestNoop(t *testing.T) {
	n := &Noop{
		ID:   "test-module-noop",
		Name: "test-module-noop",
	}
	if err := n.Validate(); err != nil {
		t.Error(err)
	}
	if err := n.Plan(); err != nil {
		t.Error(err)
	}
	if err := n.Apply(); err != nil {
		t.Error(err)
	}
	if err := n.Finalise(); err != nil {
		t.Error(err)
	}
}
