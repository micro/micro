package client

import (
	"os"
	"testing"
)

func TestBasicCall(t *testing.T) {
	if v := os.Getenv("IN_TRAVIS"); v == "yes" {
		return
	}

	response := map[string]interface{}{}
	if err := NewClient(&Options{
		Token: os.Getenv("TOKEN"),
	}).Call("helloworld", "call", map[string]interface{}{
		"name": "Alice",
	}, &response); err != nil {
		t.Fatal(err)
	}
	if len(response) > 0 {
		t.Fatal(len(response))
	}
}
