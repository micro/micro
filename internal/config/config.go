// Package config contains helper methods for
// client side config management (`~/.micro/config.json` file).
// It uses the `JSONValues` helper
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/juju/fslock"
	"github.com/micro/micro/v3/internal/user"
	conf "github.com/micro/micro/v3/service/config"
)

var (

	// lock in single process
	mtx sync.Mutex

	// File is the filepath to the config file that
	// contains all user configuration
	File = filepath.Join(user.Dir, "config.json")

	// a global lock for the config
	lock = fslock.New(File)
)

// SetConfig sets the config file
func SetConfig(configFilePath string) {
	mtx.Lock()
	defer mtx.Unlock()

	File = configFilePath
	// new lock for the file
	lock = fslock.New(File)
}

// config is a singleton which is required to ensure
// each function call doesn't load the .micro file
// from disk

// Get a value from the .micro file
func Get(path string) (string, error) {
	mtx.Lock()
	defer mtx.Unlock()

	config, err := newConfig()
	if err != nil {
		return "", err
	}

	val := config.Get(path)
	v := strings.TrimSpace(val.String(""))
	if len(v) > 0 {
		return v, nil
	}

	// try as bytes
	v = string(val.Bytes())
	v = strings.TrimSpace(v)

	// don't return nil decoded value
	if strings.TrimSpace(v) == "null" {
		return "", nil
	}

	return v, nil
}

// Set a value in the .micro file
func Set(path, value string) error {
	mtx.Lock()
	defer mtx.Unlock()

	config, err := newConfig()
	if err != nil {
		return err
	}
	// acquire lock
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	// set the value
	config.Set(path, value)

	// write to the file
	return ioutil.WriteFile(File, config.Bytes(), 0644)
}

func moveConfig(from, to string) error {
	// read the config
	b, err := ioutil.ReadFile(from)
	if err != nil {
		return fmt.Errorf("Failed to read config file %s: %v", from, err)
	}
	// remove the file
	os.Remove(from)

	// create new directory
	dir := filepath.Dir(to)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Failed to create dir %s: %v", dir, err)
	}
	// write the file to new location
	return ioutil.WriteFile(to, b, 0644)
}

// newConfig returns a loaded config
func newConfig() (*conf.JSONValues, error) {
	// check if the directory exists, otherwise create it
	dir := filepath.Dir(File)

	// for legacy purposes check if .micro is a file or directory
	if f, err := os.Stat(dir); err != nil {
		// check the error to see if the directory exists
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("Failed to create dir %s: %v", dir, err)
			}
		} else {
			return nil, fmt.Errorf("Failed to create config dir %s: %v", dir, err)
		}
	} else {
		// if not a directory, copy and move the config
		if !f.IsDir() {
			if err := moveConfig(dir, File); err != nil {
				return nil, fmt.Errorf("Failed to move config from %s to %s: %v", dir, File, err)
			}
		}
	}

	// now write the file if it does not exist
	if _, err := os.Stat(File); os.IsNotExist(err) {
		ioutil.WriteFile(File, []byte(`{"env":"local"}`), 0644)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to write config file %s: %v", File, err)
	}

	contents, err := ioutil.ReadFile(File)
	if err != nil {
		return nil, err
	}

	c := conf.NewJSONValues(contents)

	// return the conf
	return c, nil
}

func Path(paths ...string) string {
	return strings.Join(paths, ".")
}
