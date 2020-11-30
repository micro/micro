package main

import (
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/micro/micro/v3/internal/openapi/converter"
	"github.com/micro/micro/v3/service/logger"
)

func main() {

	// Get a converter:
	protoConverter := converter.New()

	// Convert the generator request:
	var ok = true
	logger.Debugf("Processing code generator request")
	res, err := protoConverter.ConvertFrom(os.Stdin)
	if err != nil {
		ok = false
		if res == nil {
			message := fmt.Sprintf("Failed to read input: %v", err)
			res = &plugin.CodeGeneratorResponse{
				Error: &message,
			}
		}
	}

	logger.Debug("Serializing code generator response")
	data, err := proto.Marshal(res)
	if err != nil {
		logger.Fatalf("Cannot marshal response: %v", err)
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		logger.Fatalf("Failed to write response: %v", err)
	}

	if ok {
		logger.Debug("Succeeded to process code generator request")
	} else {
		logger.Warn("Failed to process code generator but successfully sent the error to protoc")
		os.Exit(1)
	}
}
