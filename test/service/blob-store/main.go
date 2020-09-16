package main

import (
	"bytes"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

func main() {
	service.New(
		service.Name("blob-store"),
		service.Version("latest"),
	)

	key := uuid.New().String()

	buf := bytes.NewBuffer([]byte("world"))
	if err := store.DefaultBlobStore.Write(key, buf); err != nil {
		logger.Fatalf("Error writing to blob store: %v", err)
	} else {
		logger.Infof("Write okay")
	}

	res, err := store.DefaultBlobStore.Read(key)
	if err != nil {
		logger.Fatalf("Error reading from the blog store: %v", err)
	}
	bytes, err := ioutil.ReadAll(res)
	if err != nil {
		logger.Fatalf("Error reading result: %v", err)
	}
	logger.Infof("Read from blob store: %v", string(bytes))

	if err := store.DefaultBlobStore.Delete(key); err != nil {
		logger.Fatalf("Error deleting from blob store: %v", err)
	}
}
