package runtime

import (
	"fmt"
	"testing"
	"time"
)

func TestFmtDuration(t *testing.T) {
	tcs := []struct {
		seconds  int64
		expected string
	}{
		{seconds: 15, expected: "15s"},
		{seconds: 0, expected: "0s"},
		{seconds: 60, expected: "1m0s"},
		{seconds: 75, expected: "1m15s"},
		{seconds: 903, expected: "15m3s"},
		{seconds: 4532, expected: "1h15m32s"},
		{seconds: 82808, expected: "23h0m8s"},
		{seconds: 86400, expected: "1d0h0m0s"},
		{seconds: 1006400, expected: "11d15h33m20s"},
		{seconds: 111006360, expected: "1284d19h6m0s"},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			res := fmtDuration(time.Duration(tc.seconds) * time.Second)
			if res != tc.expected {
				t.Errorf("Expected %s but got %s", tc.expected, res)
			}
		})
	}
}
