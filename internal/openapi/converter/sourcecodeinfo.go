package converter

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Protobuf tag values for relevant message fields. Full list here:
//   https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto
const (
	tagFileDescriptorMessageType int32 = 4
	tagFileDescriptorEnumType    int32 = 5
	tagDescriptorField           int32 = 2
	tagDescriptorNestedType      int32 = 3
	tagDescriptorEnumType        int32 = 4
	tagDescriptorOneofDecl       int32 = 8
	tagEnumDescriptorValue       int32 = 2
)

type sourceCodeInfo struct {
	lookup map[proto.Message]*descriptor.SourceCodeInfo_Location
}

func (s sourceCodeInfo) GetMessage(message *descriptor.DescriptorProto) *descriptor.SourceCodeInfo_Location {
	return s.lookup[message]
}

func (s sourceCodeInfo) GetService(service *descriptor.ServiceDescriptorProto) *descriptor.SourceCodeInfo_Location {
	return s.lookup[service]
}

func (s sourceCodeInfo) GetField(field *descriptor.FieldDescriptorProto) *descriptor.SourceCodeInfo_Location {
	return s.lookup[field]
}

func (s sourceCodeInfo) GetEnum(enum *descriptor.EnumDescriptorProto) *descriptor.SourceCodeInfo_Location {
	return s.lookup[enum]
}

func (s sourceCodeInfo) GetEnumValue(value *descriptor.EnumValueDescriptorProto) *descriptor.SourceCodeInfo_Location {
	return s.lookup[value]
}

func newSourceCodeInfo(fs []*descriptor.FileDescriptorProto) *sourceCodeInfo {
	// For each source location in the provided files
	// - resolve the (annoyingly) encoded path to its message/field/service/enum/etc definition
	// - store the source info by its resolved definition
	lookup := map[proto.Message]*descriptor.SourceCodeInfo_Location{}
	for _, f := range fs {
		for _, loc := range f.GetSourceCodeInfo().GetLocation() {
			declaration := getDefinitionAtPath(f, loc.Path)
			if declaration != nil {
				lookup[declaration] = loc
			}
		}
	}
	return &sourceCodeInfo{lookup}
}

// Resolve a protobuf "file-source path" to its associated definition (eg message/field/enum/etc).
// Note that some paths don't point to definitions (some reference subcomponents like name, type,
// field #, etc) and will therefore return nil.
func getDefinitionAtPath(file *descriptor.FileDescriptorProto, path []int32) proto.Message {
	// The way protobuf encodes "file-source path" is a little opaque/tricky;
	// this doc describes how it works:
	//   https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L730

	// Starting at the root of the file descriptor, traverse its object graph by following the
	// specified path (and updating our position/state at each step) until either:
	// - we reach the definition referenced by the path (and return it)
	// - we hit a dead end because the path references a grammar element more granular than a
	//   definition (so we return nil)
	var pos proto.Message = file
	for step := 0; step < len(path); step++ {
		switch p := pos.(type) {
		case *descriptor.FileDescriptorProto:
			switch path[step] {
			case tagFileDescriptorMessageType:
				step++
				pos = p.MessageType[path[step]]
			case tagFileDescriptorEnumType:
				step++
				pos = p.EnumType[path[step]]
			default:
				return nil // ignore all other types
			}

		case *descriptor.DescriptorProto:
			switch path[step] {
			case tagDescriptorField:
				step++
				pos = p.Field[path[step]]
			case tagDescriptorNestedType:
				step++
				pos = p.NestedType[path[step]]
			case tagDescriptorEnumType:
				step++
				pos = p.EnumType[path[step]]
			case tagDescriptorOneofDecl:
				step++
				pos = p.OneofDecl[path[step]]
			default:
				return nil // ignore all other types
			}

		case *descriptor.EnumDescriptorProto:
			switch path[step] {
			case tagEnumDescriptorValue:
				step++
				pos = p.Value[path[step]]
			default:
				return nil // ignore all other types
			}

		default:
			return nil // ignore all other types
		}
	}
	return pos
}
