package converter

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/stretchr/testify/assert"

	"github.com/micro/micro/v3/internal/openapi/converter/testdata"
)

const (
	sampleProtoDirectory = "testdata/proto"
)

type sampleProto struct {
	ExpectedAPISpec  []string
	FilesToGenerate  []string
	ProtoFileName    string
	microServiceName string
}

func TestGenerateOpenAPI(t *testing.T) {

	// Configure the list of sample protos to test, and their expected API spec:
	sampleProtos := configureSampleProtos()

	// Convert the protos, compare the results against the expected API spec:
	for _, sampleProto := range sampleProtos {
		testConvertSampleProto(t, sampleProto)
	}
}

func testConvertSampleProto(t *testing.T, sampleProto sampleProto) {

	protoConverter := New()
	protoConverter.defaultSpec()

	// Open the sample proto file:
	sampleProtoFileName := fmt.Sprintf("%v/%v", sampleProtoDirectory, sampleProto.ProtoFileName)
	fileDescriptorSet := mustReadProtoFiles(t, sampleProtoDirectory, sampleProto.ProtoFileName)

	// Prepare a request:
	codeGeneratorRequest := plugin.CodeGeneratorRequest{
		FileToGenerate: sampleProto.FilesToGenerate,
		ProtoFile:      fileDescriptorSet.GetFile(),
	}

	// Perform the conversion:
	response, err := protoConverter.convert(&codeGeneratorRequest)
	assert.NoError(t, err, "Unable to convert sample proto file (%v)", sampleProtoFileName)
	assert.Equal(t, len(sampleProto.ExpectedAPISpec), len(response.File), "Incorrect number of JSON-Schema files returned for sample proto file (%v)", sampleProtoFileName)
	if len(sampleProto.ExpectedAPISpec) != len(response.File) {
		t.Fail()
	} else {
		for responseFileIndex, responseFile := range response.File {
			assert.Equal(t, strings.TrimSpace(sampleProto.ExpectedAPISpec[responseFileIndex]), *responseFile.Content, "Incorrect JSON-Schema returned for sample proto file (%v)", sampleProtoFileName)
		}
	}

	// Return now if we have no files:
	if len(response.File) == 0 {
		return
	}

	// Check for the correct prefix:
	assert.NotContains(t, response.File[0].GetName(), "samples")
}

func configureSampleProtos() map[string]sampleProto {
	return map[string]sampleProto{
		"Signup": {
			ExpectedAPISpec: []string{testdata.Signup},
			FilesToGenerate: []string{"signup.proto"},
			ProtoFileName:   "signup.proto",
		},
	}
}

// Load the specified .proto files into a FileDescriptorSet. Any errors in loading/parsing will
// immediately fail the test.
func mustReadProtoFiles(t *testing.T, includePath string, filenames ...string) *descriptor.FileDescriptorSet {
	protocBinary, err := exec.LookPath("protoc")
	if err != nil {
		t.Fatalf("Can't find 'protoc' binary in $PATH: %s", err.Error())
	}

	// Use protoc to output descriptor info for the specified .proto files.
	var args []string
	args = append(args, "--descriptor_set_out=/dev/stdout")
	args = append(args, "--include_source_info")
	args = append(args, "--include_imports")
	args = append(args, "--proto_path="+includePath)
	args = append(args, filenames...)
	cmd := exec.Command(protocBinary, args...)
	stdoutBuf := bytes.Buffer{}
	stderrBuf := bytes.Buffer{}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to load descriptor set (%s): %s: %s",
			strings.Join(cmd.Args, " "), err.Error(), stderrBuf.String())
	}
	fds := &descriptor.FileDescriptorSet{}
	err = proto.Unmarshal(stdoutBuf.Bytes(), fds)
	if err != nil {
		t.Fatalf("failed to parse protoc output as FileDescriptorSet: %s", err.Error())
	}
	return fds
}
