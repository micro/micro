package main

import (
	"bytes"
	"io/ioutil"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

func main() {
	srv := service.New(
		service.Name("blob-store"),
		service.Version("latest"),
	)

	buf := bytes.NewBuffer([]byte("world"))
	if err := store.DefaultBlobStore.Write("hello", buf); err != nil {
		logger.Fatalf("Error writing to blob store: %v", err)
	}

	res, err := store.DefaultBlobStore.Read("hello")
	if err != nil {
		logger.Fatalf("Error reading from the blog store: %v", err)
	}
	bytes, err := ioutil.ReadAll(res)
	if err != nil {
		logger.Fatalf("Error reading result: %v", err)
	}
	logger.Infof("Read from blob store: %v", string(bytes))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
