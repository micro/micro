// +build integration

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	TrySuite(t, testConfig, retryCount)
}

func testConfig(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	t.T().Run("Test non existing record", func(tee *testing.T) {
		if err := Try("Calling micro config read", t, func() ([]byte, error) {
			outp, err := cmd.Exec("config", "get", "somekey")
			if err == nil {
				return outp, errors.New("config gete should fail")
			}
			if !strings.Contains(string(outp), "Not found") {
				return outp, fmt.Errorf("Output should be 'Not found\n', got %v", string(outp))
			}
			return outp, nil
		}, 5*time.Second); err != nil {
			return
		}
	})

	t.T().Run("Test no dot get set delete", func(tee *testing.T) {
		// This needs to be retried to the the "error listing rules"
		// error log output that happens when the auth service is not yet available.

		if err := Try("Calling micro config set", t, func() ([]byte, error) {
			outp, err := cmd.Exec("config", "set", "somekey", "val1")
			if err != nil {
				return outp, err
			}
			if string(outp) != "" {
				return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
			}
			return outp, err
		}, 5*time.Second); err != nil {
			return
		}

		if err := Try("micro config get somekey", t, func() ([]byte, error) {
			outp, err := cmd.Exec("config", "get", "somekey")
			if err != nil {
				return outp, err
			}
			if string(outp) != "val1\n" {
				return outp, errors.New("Expected 'val1\n'")
			}
			return outp, err
		}, 8*time.Second); err != nil {
			return
		}

		outp, err := cmd.Exec("config", "del", "somekey")
		if err != nil {
			t.Fatalf(string(outp))
			return
		}
		if string(outp) != "" {
			t.Fatalf("Expected '', got: '%v'", string(outp))
			return
		}

		if err := Try("micro config get somekey", t, func() ([]byte, error) {
			outp, err := cmd.Exec("config", "get", "somekey")
			if err == nil {
				return outp, errors.New("getting somekey should fail")
			}
			if !strings.Contains(string(outp), "Not found") {
				return outp, errors.New("Expected 'Not found'")
			}
			return outp, nil
		}, 8*time.Second); err != nil {
			return
		}

		t.T().Run("Test dot notation", func(tee *testing.T) {
			outp, err = cmd.Exec("config", "set", "someotherkey.subkey", "otherval1")
			if err != nil {
				t.Fatal(err)
				return
			}
			if string(outp) != "" {
				t.Fatalf("Expected no output, got: %v", string(outp))
				return
			}

			if err := Try("micro config get someotherkey.subkey", t, func() ([]byte, error) {
				outp, err := cmd.Exec("config", "get", "someotherkey.subkey")
				if err != nil {
					return outp, err
				}
				if string(outp) != "otherval1\n" {
					return outp, errors.New("Expected 'otherval1\n'")
				}
				return outp, err
			}, 8*time.Second); err != nil {
				return
			}
		})
	})

	t.T().Run("Test basic merge", func(tee *testing.T) {
		outp, err := cmd.Exec("config", "set", "mergekey", `{"a":1}`)
		if err != nil {
			t.Fatal(err)
			return
		}

		if string(outp) != "" {
			t.Fatalf("Expected no output, got: %v", string(outp))
			return
		}

		outp, err = cmd.Exec("config", "set", "mergekey", `{"b":2}`)
		if err != nil {
			t.Fatal(err)
			return
		}
		if string(outp) != "" {
			t.Fatalf("Expected no output, got: %v", string(outp))
			return
		}

		outp, err = cmd.Exec("config", "get", "mergekey")
		if err != nil {
			t.Fatal(err)
			return
		}
		var m map[string]interface{}
		err = json.Unmarshal(outp, &m)
		if err != nil {
			t.Fatal(err)
			return
		}
		expected := map[string]interface{}{
			"a": float64(1),
			"b": float64(2),
		}
		if !reflect.DeepEqual(m, expected) {
			t.Fatalf("Output is %v, expected: %v", m, expected)
			return
		}
	})

	t.T().Run("Test config escape", func(tee *testing.T) {
		outp, err := cmd.Exec("config", "set", "jsonescape", `"Value with <> signs"`)
		if err != nil {
			t.Fatal(err)
			return
		}
		if string(outp) != "" {
			t.Fatalf("Expected no output, got: %v", string(outp))
			return
		}
		outp, err = cmd.Exec("config", "get", "jsonescape")
		if err != nil {
			t.Fatal(err)
			return
		}

		if strings.TrimSpace(string(outp)) != "Value with <> signs" {
			t.Fatalf("Expected 'Value with <> signs', got: '%v'", string(outp))
			return
		}
	})

	t.T().Run("Test complex merge", func(tee *testing.T) {
		outp, err := cmd.Exec("config", "set", "complexmerge", `{"a":1,"b":{"b1":2},"c":3}`)
		if err != nil {
			t.Fatal(err)
			return
		}
		if string(outp) != "" {
			t.Fatalf("Expected no output, got: %v", string(outp))
			return
		}
		outp, err = cmd.Exec("config", "set", "complexmerge", `{"d":4,"b":{"b2":2.2}}`)
		if err != nil {
			t.Fatal(err)
			return
		}

		expected := map[string]interface{}{
			"a": float64(1),
			"b": map[string]interface{}{
				"b1": float64(2),
				"b2": float64(2.2),
			},
			"c": float64(3),
			"d": float64(4),
		}
		outp, err = cmd.Exec("config", "get", "complexmerge")
		if err != nil {
			t.Fatal(err)
			return
		}
		result := map[string]interface{}{}
		err = json.Unmarshal(outp, &result)
		if err != nil {
			t.Fatal(err)
			return
		}
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("Expected %v is different from result %v", expected, result)
			return
		}
	})

	t.T().Run("Test plain old type being overwritten by map", func(tee *testing.T) {
		outp, err := cmd.Exec("config", "set", "pod", "hi")
		if err != nil {
			t.Fatal(err)
			return
		}
		if string(outp) != "" {
			t.Fatalf("Expected no output, got: %v", string(outp))
			return
		}
		outp, err = cmd.Exec("config", "set", "pod", `{"a":1}`)
		if err != nil {
			t.Fatal(err)
			return
		}

		expected := map[string]interface{}{
			"a": float64(1),
		}
		outp, err = cmd.Exec("config", "get", "pod")
		if err != nil {
			t.Fatal(err)
			return
		}
		result := map[string]interface{}{}
		err = json.Unmarshal(outp, &result)
		if err != nil {
			t.Fatal(err)
			return
		}
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("Expected %v is different from result %v", expected, result)
			return
		}
	})

	t.T().Run("Test plain old type being overwritten by pod", func(tee *testing.T) {
		outp, err := cmd.Exec("config", "set", "pod", "hi")
		if err != nil {
			t.Fatal(err)
			return
		}
		if string(outp) != "" {
			t.Fatalf("Expected no output, got: %v", string(outp))
			return
		}
		outp, err = cmd.Exec("config", "set", "pod", "hello")
		if err != nil {
			tee.Fatal(err)
			return
		}
		outp, err = cmd.Exec("config", "get", "pod")
		if err != nil {
			t.Fatal(err)
			return
		}
		expected := "hello"
		if !reflect.DeepEqual(strings.TrimSpace(string(outp)), expected) {
			tee.Fatalf("Expected %v is different from result %v", expected, strings.TrimSpace(string(outp)))
			return
		}
	})
}

func TestConfigReadFromService(t *testing.T) {
	TrySuite(t, testConfigReadFromService, retryCount)
}

func testConfigReadFromService(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// This needs to be retried to the the "error listing rules"
	// error log output that happens when the auth service is not yet available.
	if err := Try("Calling micro config set", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "set", "key.subkey", "val1")
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		outp, err = cmd.Exec("config", "set", "--secret", "key.subkey1", "42")
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		// Testing JSON escape of "<" etc chars.
		outp, err = cmd.Exec("config", "set", "--secret", "key.subkey2", "\"Micro Team <support@m3o.com>\"")
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		// Setting an other key for `val.String` test
		outp, err = cmd.Exec("config", "set", "key.subkey3", "\"Micro Test <test@m3o.com>\"")
		if err != nil {
			return outp, err
		}
		if string(outp) != "" {
			return outp, fmt.Errorf("Expected no output, got: %v", string(outp))
		}
		return outp, err
	}, 5*time.Second); err != nil {
		return
	}

	// check the value set correctly
	if err := Try("Calling micro config get", t, func() ([]byte, error) {
		outp, err := cmd.Exec("config", "get", "key")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "val1") {
			return outp, fmt.Errorf("Expected output to contain val1, got: %v", string(outp))
		}

		return outp, err
	}, 5*time.Second); err != nil {
		return
	}

	outp, err := cmd.Exec("run", "./services/test/conf")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Try logs read", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "conf")
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "val1") {
			return outp, fmt.Errorf("Expected val1 in output, got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "42") {
			return outp, fmt.Errorf("Expected output to contain 42, got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "Micro Team <support@m3o.com>") {
			return outp, fmt.Errorf("Expected output to contain \"Micro Team <support@m3o.com>\", got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "Micro Test <test@m3o.com>") {
			return outp, fmt.Errorf("Expected output to contain \"Micro Test <test@m3o.com>\", got: %v", string(outp))
		}
		if !strings.Contains(string(outp), "Default Hello") {
			return outp, fmt.Errorf("Expected output to contain \"Default Hello\", got: %v", string(outp))
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}
}
