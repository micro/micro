package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/juju/fslock"
	conf "github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/config/source/file"
)

var (
	// FileName for global micro config
	FileName = ".micro/config.json"

	path, _ = filePath()

	// a global lock for the config
	lock = fslock.New(path)
)

// config is a singleton which is required to ensure
// each function call doesn't load the .micro file
// from disk

// Get a value from the .micro file
func Get(path ...string) (string, error) {
	config, err := newConfig()
	if err != nil {
		return "", err
	}

	// acquire lock
	if err := lock.LockWithTimeout(time.Second); err != nil {
		return "", err
	}
	defer lock.Unlock()

	val := config.Get(path...)
	v := strings.TrimSpace(val.String(""))
	if len(v) > 0 {
		return v, nil
	}

	// try as bytes
	v = string(val.Bytes())
	v = strings.TrimSpace(v)

	// don't return nil decoded value
	if v == "null" {
		return "", nil
	}

	return v, nil
}

// Set a value in the .micro file
func Set(value string, path ...string) error {
	// get the filepath
	fp, err := filePath()
	if err != nil {
		return err
	}

	config, err := newConfig()
	if err != nil {
		return err
	}

	// acquire lock
	if err := lock.LockWithTimeout(time.Second); err != nil {
		return err
	}
	defer lock.Unlock()

	// set the value
	config.Set(value, path...)

	// write to the file
	return ioutil.WriteFile(fp, config.Bytes(), 0644)
}

func filePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, FileName), nil
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
func newConfig() (conf.Config, error) {
	// get the filepath
	fp, err := filePath()
	if err != nil {
		return nil, err
	}

	// check if the directory exists, otherwise create it
	dir := filepath.Dir(fp)

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
			if err := moveConfig(dir, fp); err != nil {
				return nil, fmt.Errorf("Failed to move config from %s to %s: %v", dir, fp, err)
			}
		}
	}

	// now write the file if it does not exist
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		ioutil.WriteFile(fp, []byte(`{}`), 0644)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to write config file %s: %v", fp, err)
	}

	// create a new config
	c, err := conf.NewConfig(
		conf.WithSource(
			file.NewSource(
				file.WithPath(fp),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// load the config
	if err := c.Load(); err != nil {
		return nil, err
	}

	// return the conf
	return c, nil
}
