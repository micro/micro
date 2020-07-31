package config

import (
	"os"
	"testing"

	"github.com/juju/fslock"
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

	// change the config path for the lifetime
	// of this test
	savePath := FileName
	saveLock := lock
	FileName = ".micro/config-test.json"
	path, _ := filePath()
	lock = fslock.New(path)

	defer func() {
		FileName = savePath
		lock = saveLock
	}()

	fp, err := filePath()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer os.Remove(fp)

			if _, err := os.Stat(fp); err != os.ErrNotExist {
				os.Remove(fp)
			}

			for k, v := range tc.values {
				if err := Set(v, k); err != nil {
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
