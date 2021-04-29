package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaPaths(t *testing.T) {

	// URL path (including service-name, proto-service, and method):
	assert.Equal(t, "/tests/schemaPaths/urlPath", urlPath("tests", "schemaPaths", "urlPath"))

	// Service name (last word in a dot-delimited string):
	assert.Equal(t, "schemapaths", protoServiceName("tests.schemapaths"))

	// How to find a message schema in an OpenAPI spec:
	assert.Equal(t, "#/components/schemas/urlPathPayload", messageSchemaPath("urlPathPayload"))

	// The name of a request body schema:
	assert.Equal(t, "schemaPathsurlPathRequest", requestBodyName("schemaPaths", "urlPath"))

	// How to find a request body schema in an OpenAPI spec:
	assert.Equal(t, "#/components/requestBodies/schemaPathsurlPathRequest", requestBodySchemaPath("schemaPathsurlPathRequest"))

	// The name of a response body schema:
	assert.Equal(t, "schemaPathsurlPathResponse", responseBodyName("schemaPaths", "urlPath"))

	// How to find a response body schema in an OpenAPI spec:
	assert.Equal(t, "#/components/responses/schemaPathsurlPathResponse", responseBodySchemaPath("schemaPathsurlPathResponse"))
}
