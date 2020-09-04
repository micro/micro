package user

import (
	"log"
	"os/user"
	"path/filepath"
)

var (
	Dir  = ""
	path = ".micro"
)

func init() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	Dir = filepath.Join(user.HomeDir, path)
}
