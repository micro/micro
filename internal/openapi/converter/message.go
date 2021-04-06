package converter

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/micro/micro/v3/service/logger"
)

const (
	openAPIFormatByte     = "byte"
	openAPIFormatDateTime = "date-time"
	openAPIFormatDouble   = "double"
	openAPIFormatInt32    = "int32"
	openAPIFormatInt64    = "int64"
	openAPITypeArray      = "array"
	openAPITypeBoolean    = "boolean"
	openAPITypeNumber     = "number"
	openAPITypeObject     = "object"
	openAPITypeString     = "string"
)

var (
	globalPkg = &ProtoPackage{
		name:     "",
		parent:   nil,
		children: make(map[string]*ProtoPackage),
		types:    make(map[string]*descriptor.DescriptorProto),
	}

	wellKnownTypes = map[string]bool{
		"DoubleValue": true,
		"FloatValue":  true,
		"Int64Value":  true,
		"UInt64Value": true,
		"Int32Value":  true,
		"UInt32Value": true,
		"BoolValue":   true,
		"StringValue": true,
		"BytesValue":  true,
		"Value":       true,
	}
)

func (c *Converter) registerType(pkgName *string, msg *descriptor.DescriptorProto) {
	pkg := globalPkg
	if pkgName != nil {
		for _, node := range strings.Split(*pkgName, ".") {
			if pkg == globalPkg && node == "" {
				// Skips leading "."
				continue
			}
			child, ok := pkg.children[node]
			if !ok {
				child = &ProtoPackage{
					name:     pkg.name + "." + node,
					parent:   pkg,
					children: make(map[string]*ProtoPackage),
					types:    make(map[string]*descriptor.DescriptorProto),
				}
				pkg.children[node] = child
			}
			pkg = child
		}
	}
	pkg.types[msg.GetName()] = msg
}

func (c *Converter) relativelyLookupNestedType(desc *descriptor.DescriptorProto, name string) (*descriptor.DescriptorProto, bool) {
	components := strings.Split(name, ".")
componentLoop:
	for _, component := range components {
		for _, nested := range desc.GetNestedType() {
			if nested.GetName() == component {
				desc = nested
				continue componentLoop
			}
		}
		logger.Infof("no such nested message (%s.%s)", component, desc.GetName())
		return nil, false
	}
	return desc, true
}

// Convert a proto "field" (essentially a type-switch with some recursion):
func (c *Converter) convertField(curPkg *ProtoPackage, desc *descriptor.FieldDescriptorProto, msg *descriptor.DescriptorProto) (*openapi3.Schema, error) {

	// Prepare a new jsonschema.Type for our eventual return value:
	componentSchema := &openapi3.Schema{}

	// Generate a description from src comments (if available)
	if src := c.sourceInfo.GetField(desc); src != nil {
		componentSchema.Description = formatDescription(src)
	}

	// Switch the types, and pick a JSONSchema equivalent:
	switch desc.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		componentSchema.Type = openAPITypeNumber
		componentSchema.Format = openAPIFormatDouble

	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SINT32:
		componentSchema.Type = openAPITypeNumber
		componentSchema.Format = openAPIFormatInt32

	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		componentSchema.Type = openAPITypeNumber
		componentSchema.Format = openAPIFormatInt64

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		componentSchema.Type = openAPITypeString

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		componentSchema.Type = openAPITypeString
		componentSchema.Format = openAPIFormatByte

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		componentSchema.Type = "string"

		// Go through all the enums we have, see if we can match any to this field by name:
		for _, enumDescriptor := range msg.GetEnumType() {

			// Each one has several values:
			for _, enumValue := range enumDescriptor.Value {

				// Figure out the entire name of this field:
				fullFieldName := fmt.Sprintf(".%v.%v", *msg.Name, *enumDescriptor.Name)

				// If we find ENUM values for this field then put them into the JSONSchema list of allowed ENUM values:
				if strings.HasSuffix(desc.GetTypeName(), fullFieldName) {
					componentSchema.Enum = append(componentSchema.Enum, *enumValue.Name)
				}
			}
		}

	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		componentSchema.Type = openAPITypeBoolean

	case descriptor.FieldDescriptorProto_TYPE_GROUP, descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		switch desc.GetTypeName() {
		case ".google.protobuf.Timestamp":
			componentSchema.Type = openAPITypeString
			componentSchema.Format = openAPIFormatDateTime
		default:
			componentSchema.Type = openAPITypeObject
		}

	default:
		return nil, fmt.Errorf("unrecognized field type: %s", desc.GetType().String())
	}

	// Recurse array of primitive types:
	if desc.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED && componentSchema.Type != openAPITypeObject {
		componentSchema.Items = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: componentSchema.Type,
			},
		}
		componentSchema.Type = openAPITypeArray
		return componentSchema, nil
	}

	// Recurse nested objects / arrays of objects (if necessary):
	if componentSchema.Type == openAPITypeObject {

		recordType, pkgName, ok := c.lookupType(curPkg, desc.GetTypeName())
		if !ok {
			return nil, fmt.Errorf("no such message type named %s", desc.GetTypeName())
		}

		// Recurse the recordType:
		recursedComponentSchema, err := c.recursiveConvertMessageType(curPkg, recordType, pkgName)
		if err != nil {
			return nil, err
		}

		// Maps, arrays, and objects are structured in different ways:
		switch {
		// Arrays:
		case desc.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED:
			componentSchema.Items = &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       openAPITypeObject,
					Properties: recursedComponentSchema.Properties,
				},
			}
			componentSchema.Type = openAPITypeArray
		// Maps:
		case recordType.Options.GetMapEntry():
			logger.Tracef("Found a map (%s.%s)", *msg.Name, recordType.GetName())
			componentSchema.Type = openAPITypeObject
			componentSchema.AdditionalProperties = openapi3.NewSchemaRef("", recursedComponentSchema)

		// Objects:
		default:
			componentSchema.Properties = recursedComponentSchema.Properties
			// recursedComponentSchemaRef := fmt.Sprintf("#/components/schemas/%s", recursedComponentSchema.Title)
			// componentSchema.Properties = openapi3.NewSchemaRef(recursedComponentSchemaRef, nil)
		}
	}

	return componentSchema, nil
}

// Converts a proto "MESSAGE" into an OpenAPI schema:
func (c *Converter) convertMessageType(curPkg *ProtoPackage, msg *descriptor.DescriptorProto) (*openapi3.Schema, error) {

	// main schema for the message
	rootType, err := c.recursiveConvertMessageType(curPkg, msg, "")
	if err != nil {
		return nil, err
	}

	return rootType, nil
}

type nameAndCounter struct {
	name    string
	counter int
}

func (c *Converter) recursiveConvertMessageType(curPkg *ProtoPackage, msg *descriptor.DescriptorProto, pkgName string) (*openapi3.Schema, error) {
	if msg.Name != nil && wellKnownTypes[*msg.Name] && pkgName == ".google.protobuf" {
		componentSchema := &openapi3.Schema{
			Title: msg.GetName(),
		}
		switch *msg.Name {
		case "DoubleValue", "FloatValue":
			componentSchema.Type = openAPITypeNumber
			componentSchema.Format = openAPIFormatDouble
		case "Int32Value", "UInt32Value":
			componentSchema.Type = openAPITypeNumber
			componentSchema.Format = openAPIFormatInt32
		case "Int64Value", "UInt64Value":
			componentSchema.Type = openAPITypeNumber
			componentSchema.Format = openAPIFormatInt64
		case "BoolValue":
			componentSchema.Type = openAPITypeBoolean
		case "StringValue":
			componentSchema.Type = openAPITypeString
		case "BytesValue":
			componentSchema.Type = openAPITypeString
			componentSchema.Format = openAPIFormatByte
		case "Value":
			componentSchema.Type = openAPITypeObject
		}
		return componentSchema, nil
	}

	// Prepare a new jsonschema:
	componentSchema := &openapi3.Schema{
		Properties: make(map[string]*openapi3.SchemaRef),
		Title:      msg.GetName(),
		Type:       openAPITypeObject,
	}

	// Generate a description from src comments (if available)
	if src := c.sourceInfo.GetMessage(msg); src != nil {
		componentSchema.Description = formatDescription(src)
	}

	logger.Tracef("Converting message (%s)", proto.MarshalTextString(msg))

	// Recurse each field:
	for _, fieldDesc := range msg.GetField() {
		recursedComponentSchema, err := c.convertField(curPkg, fieldDesc, msg)
		if err != nil {
			logger.Errorf("Failed to convert field (%s.%s): %v", msg.GetName(), fieldDesc.GetName(), err)
			return nil, err
		}
		logger.Tracef("Converted field: %s => %s", fieldDesc.GetName(), recursedComponentSchema.Type)

		// Add it to the properties (by its JSON name):
		componentSchema.Properties[fieldDesc.GetJsonName()] = openapi3.NewSchemaRef("", recursedComponentSchema)
	}

	return componentSchema, nil
}

func formatDescription(sl *descriptor.SourceCodeInfo_Location) string {
	var lines []string
	for _, str := range sl.GetLeadingDetachedComments() {
		if s := strings.TrimSpace(str); s != "" {
			lines = append(lines, s)
		}
	}
	if s := strings.TrimSpace(sl.GetLeadingComments()); s != "" {
		lines = append(lines, s)
	}
	if s := strings.TrimSpace(sl.GetTrailingComments()); s != "" {
		lines = append(lines, s)
	}
	return strings.Join(lines, "\n\n")
}
