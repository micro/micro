package config

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	tt := []struct {
		name   string
		values map[string]string
	}{
		{
			name: "No values",
		},
		{
			name: "Single value",
			values: map[string]string{
				"foo": "bar",
			},
		},
		{
			name: "Multiple values",
			values: map[string]string{
				"foo":   "bar",
				"apple": "tree",
			},
		},
	}

	fp, err := filePath()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := os.Stat(fp); err != os.ErrNotExist {
				os.Remove(fp)
			}

			for k, v := range tc.values {
				if err := Set(k, v); err != nil {
					t.Error(err)
				}
			}

			for k, v := range tc.values {
				val, err := Get(k)
				if err != nil {
					t.Error(err)
					continue
				}

				if v != val {
					t.Errorf("Got '%v' but expected '%v'", val, v)
				}
			}
		})
	}
}
