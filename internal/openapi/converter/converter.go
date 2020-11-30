package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/micro/micro/v3/service/logger"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Converter is everything you need to convert Micro protos into an OpenAPI spec:
type Converter struct {
	microServiceName string
	openAPISpec      *openapi3.Swagger
	sourceInfo       *sourceCodeInfo
}

// New returns a configured converter:
func New() *Converter {
	return &Converter{}
}

// ConvertFrom tells the convert to work on the given input:
func (c *Converter) ConvertFrom(rd io.Reader) (*plugin.CodeGeneratorResponse, error) {
	logger.Debug("Reading code generation request")
	input, err := ioutil.ReadAll(rd)
	if err != nil {
		logger.Errorf("Failed to read request: %v", err)
		return nil, err
	}

	req := &plugin.CodeGeneratorRequest{}
	err = proto.Unmarshal(input, req)
	if err != nil {
		logger.Errorf("Can't unmarshal input: %v", err)
		return nil, err
	}

	c.defaultSpec()

	logger.Debugf("Converting input: %v", err)
	return c.convert(req)
}

// Converts a proto file into an OpenAPI spec:
func (c *Converter) convertFile(file *descriptor.FileDescriptorProto) error {

	// Input filename:
	protoFileName := path.Base(file.GetName())

	// Get the proto package:
	pkg, ok := c.relativelyLookupPackage(globalPkg, file.GetPackage())
	if !ok {
		return fmt.Errorf("no such package found: %s", file.GetPackage())
	}
	c.openAPISpec.Info.Title = strings.Title(strings.Replace(pkg.name, ".", "", 1))

	// Process messages:
	for _, msg := range file.GetMessageType() {

		// Convert the message:
		logger.Debugf("Generating component schema for message (%s) from proto file (%s)", msg.GetName(), protoFileName)
		componentSchema, err := c.convertMessageType(pkg, msg)
		if err != nil {
			logger.Errorf("Failed to convert (%s): %v", protoFileName, err)
			return err
		}

		// Add the message to our component schemas (we'll refer to these later when we build the service endpoints):
		// componentSchemaKey := fmt.Sprintf("%s.%s", pkg.name, componentSchema.Title)
		c.openAPISpec.Components.Schemas[componentSchema.Title] = openapi3.NewSchemaRef("", componentSchema)
	}

	// Process services:
	for _, svc := range file.GetService() {

		// Convert the service:
		logger.Infof("Generating service (%s) from proto file (%s)", svc.GetName(), protoFileName)
		servicePaths, err := c.convertServiceType(file, pkg, svc)
		if err != nil {
			logger.Errorf("Failed to convert (%s): %v", protoFileName, err)
			return err
		}

		// Add the paths to our API:
		for path, pathItem := range servicePaths {
			c.openAPISpec.Paths[path] = pathItem
		}
	}

	return nil
}

func (c *Converter) convert(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	res := &plugin.CodeGeneratorResponse{}

	c.parseGeneratorParameters(req.GetParameter())

	// Parse the source code (this is where we pick up code comments, which become schema descriptions):
	c.sourceInfo = newSourceCodeInfo(req.GetProtoFile())

	generateTargets := make(map[string]bool)
	for _, file := range req.GetFileToGenerate() {
		generateTargets[file] = true
	}

	// We're potentially dealing with several proto files:
	for _, file := range req.GetProtoFile() {

		// Make sure it belongs to a package (sometimes they don't):
		if file.GetPackage() == "" {
			logger.Warnf("Proto file (%s) doesn't specify a package", file.GetName())
			continue
		}

		// Set the service name from the proto package (if it isn't already set):
		if c.microServiceName == "" {
			c.microServiceName = protoServiceName(file.GetPackage())
		}

		// Register all of the messages we can find:
		for _, msg := range file.GetMessageType() {
			logger.Debugf("Loading a message (%s/%s)", file.GetPackage(), msg.GetName())
			c.registerType(file.Package, msg)
		}

		if _, ok := generateTargets[file.GetName()]; ok {
			logger.Debugf("Converting file (%s)", file.GetName())
			if err := c.convertFile(file); err != nil {
				res.Error = proto.String(fmt.Sprintf("Failed to convert %s: %v", file.GetName(), err))
				return res, err
			}
		}
	}

	// Marshal the OpenAPI spec:
	marshaledSpec, err := json.MarshalIndent(c.openAPISpec, "", "  ")
	if err != nil {
		logger.Errorf("Unable to marshal the OpenAPI spec: %v", err)
		return nil, err
	}

	// Add a response file:
	res.File = []*plugin.CodeGeneratorResponse_File{
		{
			Name:    proto.String(c.openAPISpecFileName()),
			Content: proto.String(string(marshaledSpec)),
		},
	}

	return res, nil
}

func (c *Converter) openAPISpecFileName() string {
	return fmt.Sprintf("api-%s.json", c.microServiceName)
}

func (c *Converter) parseGeneratorParameters(parameters string) {
	logger.Debug("Parsing params")

	for _, parameter := range strings.Split(parameters, ",") {

		logger.Debugf("Param: %s", parameter)

		// Allow users to specify the service name:
		if serviceNameParameter := strings.Split(parameter, "service="); len(serviceNameParameter) == 2 {
			c.microServiceName = serviceNameParameter[1]
			logger.Infof("Service name: %s", c.microServiceName)
		}
	}
}
