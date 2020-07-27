package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	conf "github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/config/source/file"
	"github.com/micro/go-micro/v3/util/log"
)

// FileName for global micro config
const FileName = ".micro"

// config is a singleton which is required to ensure
// each function call doesn't load the .micro file
// from disk
var config = newConfig()

// Get a value from the .micro file
func Get(path ...string) (string, error) {
	c := newConfig()
	val := c.Get(path...)
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

	// set the value
	c := newConfig()
	c.Set(value, path...)

	// write to the file
	return ioutil.WriteFile(fp, c.Bytes(), 0644)
}

func filePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, FileName), nil
}

// newConfig returns a loaded config
func newConfig() conf.Config {
	// get the filepath
	fp, err := filePath()
	if err != nil {
		log.Error(err)
		return conf.DefaultConfig
	}

	// write the file if it does not exist
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		ioutil.WriteFile(fp, []byte{}, 0644)
	} else if err != nil {
		log.Error(err)
		return conf.DefaultConfig
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
		log.Error(err)
		return conf.DefaultConfig
	}

	// load the config
	if err := c.Load(); err != nil {
		log.Error(err)
		return conf.DefaultConfig
	}

	// return the conf
	return c
}
