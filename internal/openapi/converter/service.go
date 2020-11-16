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

		c.logger.Debugf("Processing method %s.%s()", svc.GetName(), method.GetName())

		// The URL path is the service name and method name:
		path := fmt.Sprintf("/%s/%s", svc.GetName(), method.GetName())

		// See if we can get the request paylod schema:
		requestBodySchema, ok := c.openAPISpec.Components.Schemas[*method.InputType]
		if !ok {
			c.logger.Warnf("Couldn't find request body payload (%s)", *method.InputType)
			continue
		}

		// Make a request body:
		requestBody := &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.NewContentWithJSONSchemaRef(requestBodySchema),
			},
		}

		// Add it to the spec:
		c.openAPISpec.Components.RequestBodies[*method.InputType] = requestBody

		// // See if we can get the response paylod schema:
		// responseBodySchema, ok := c.componentSchemas[*method.OutputType]
		// if !ok {
		// 	c.logger.Warnf("Couldn't find response body payload (%s)", *method.OutputType)
		// 	continue
		// }

		// Prepare a path item based on these payloads:
		pathItem := &openapi3.PathItem{
			Summary: fmt.Sprintf("%s: %s.%s()", file.GetName(), svc.GetName(), method.GetName()),
			Post: &openapi3.Operation{
				RequestBody: requestBody,
				Responses:   openapi3.Responses{},
				Security: &openapi3.SecurityRequirements{
					{
						"MicroAPIToken": []string{},
					},
				},
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
