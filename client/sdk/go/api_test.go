package api

import (
	"os"
	"testing"
)

func TestBasicCall(t *testing.T) {
	if v := os.Getenv("IN_TRAVIS_CI"); v == "yes" {
		return
	}

	response := map[string]interface{}{}
	if err := NewClient(&Options{
		Token: os.Getenv("MICRO_API_TOKEN"),
	}).Call("helloworld", "call", map[string]interface{}{
		"name": "Alice",
	}, &response); err != nil {
		t.Fatal(err)
	}
	if len(response) > 0 {
		t.Fatal(len(response))
	}
}
