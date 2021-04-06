package converter

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/micro/micro/v3/service/logger"
)

// Converts a proto "SERVICE" into an OpenAPI path:
func (c *Converter) convertServiceType(file *descriptor.FileDescriptorProto, curPkg *ProtoPackage, svc *descriptor.ServiceDescriptorProto) (map[string]*openapi3.PathItem, error) {
	pathItems := make(map[string]*openapi3.PathItem)

	// Add a path item for each method in the service:
	for _, method := range svc.GetMethod() {
		logger.Debugf("Processing method %s.%s()", svc.GetName(), method.GetName())

		// Figure out the URL path:
		path := urlPath(c.microServiceName, svc.GetName(), method.GetName())

		// We need to reformat the request name to match what is produced by the message converter:
		requestPayloadSchemaName := protoServiceName(*method.InputType)

		// See if we can get the request paylod schema:
		if _, ok := c.openAPISpec.Components.Schemas[requestPayloadSchemaName]; !ok {
			logger.Warnf("Couldn't find request body payload (%s)", requestPayloadSchemaName)
			continue
		}

		// Make a request body:
		requestBodyName := requestBodyName(svc.GetName(), method.GetName())
		requestBody := &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: messageSchemaPath(requestPayloadSchemaName),
						},
					},
				},
				Description: requestBodyName,
			},
		}

		// Add it to the spec:
		c.openAPISpec.Components.RequestBodies[requestBodyName] = requestBody

		// We need to reformat the response name to match what is produced by the message converter:
		responsePayloadSchemaName := protoServiceName(*method.OutputType)

		// See if we can get the response paylod schema:
		if _, ok := c.openAPISpec.Components.Schemas[responsePayloadSchemaName]; !ok {
			logger.Warnf("Couldn't find response body payload (%s)", responsePayloadSchemaName)
			continue
		}

		// Make a response body:
		responseBodyName := responseBodyName(svc.GetName(), method.GetName())
		responseBody := &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: messageSchemaPath(responsePayloadSchemaName),
						},
					},
				},
				Description: &responseBodyName,
			},
		}

		// Add it to the spec:
		c.openAPISpec.Components.Responses[responseBodyName] = responseBody

		// Prepare a path item based on these payloads:
		pathItem := &openapi3.PathItem{
			Parameters: openapi3.Parameters{
				{
					Value: &openapi3.Parameter{
						In:       "header",
						Name:     "Micro-Namespace",
						Required: true,
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
					},
				},
			},
			Post: &openapi3.Operation{
				RequestBody: &openapi3.RequestBodyRef{
					Ref: requestBodySchemaPath(requestBodyName),
				},
				Responses: openapi3.Responses{
					"default": &openapi3.ResponseRef{
						Ref: responseBodySchemaPath("MicroAPIError"),
					},
					"200": &openapi3.ResponseRef{
						Ref: responseBodySchemaPath(responseBodyName),
					},
				},
				Security: &openapi3.SecurityRequirements{
					{
						"MicroAPIToken": []string{},
					},
				},
				Summary: fmt.Sprintf("%s.%s(%s)", svc.GetName(), method.GetName(), requestPayloadSchemaName),
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
