package converter

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Converts a proto "SERVICE" into an OpenAPI path:
func (c *Converter) convertServiceType(file *descriptor.FileDescriptorProto, curPkg *ProtoPackage, svc *descriptor.ServiceDescriptorProto) (map[string]*openapi3.PathItem, error) {

	pathItems := make(map[string]*openapi3.PathItem)

	// Add a path item for each method in the service:
	for _, method := range svc.GetMethod() {

		path := fmt.Sprintf("/%s/%s", svc.GetName(), method.GetName())

		pathItem := &openapi3.PathItem{
			Summary: fmt.Sprintf("%s: %s.%s()", file.GetName(), svc.GetName(), method.GetName()),
			Post: &openapi3.Operation{
				RequestBody: &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: openapi3.NewContentWithJSONSchema(c.componentSchemas[*method.InputType]),
					},
				},
				Responses: openapi3.Responses{},
			},
		}

		// Generate a description from src comments (if available)
		if src := c.sourceInfo.GetService(svc); src != nil {
			pathItem.Description = formatDescription(src)
		}

		pathItems[path] = pathItem
	}

	return pathItems, nil
}
