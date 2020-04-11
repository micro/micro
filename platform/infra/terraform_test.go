package infra

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestTerraformModule(t *testing.T) {
	if e := os.Getenv("IN_TRAVIS_CI"); e == "yes" {
		t.Skip("In Travis CI")
	}
	rand.Seed(time.Now().UnixNano())
	path := fmt.Sprintf("/tmp/test-module-%d", rand.Int31())
	testModule := &TerraformModule{
		ID:     fmt.Sprintf("test-module-%d", rand.Int31()),
		Name:   "test",
		Path:   path,
		Source: "./network",
		DryRun: true,
	}
	defer testModule.Finalise()
	if err := testModule.Validate(); err != nil {
		t.Error(err)
	}
	if err := testModule.Plan(); err != nil {
		t.Error(err)
	}
	if err := testModule.Apply(); err != nil {
		t.Error(err)
	}
}
