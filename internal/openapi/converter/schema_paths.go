package converter

import (
	"fmt"
	"strings"
)

func payloadSchemaName(inputType string) string {
	return strings.Split(inputType, ".")[2]
}

func messageSchemaPath(schemaName string) string {
	return fmt.Sprintf("#/components/schemas/%s", schemaName)
}

func requestBodyName(serviceName, methodName string) string {
	return fmt.Sprintf("%s%sRequest", serviceName, methodName)
}

func requestBodySchemaPath(requestBodyName string) string {
	return fmt.Sprintf("#/components/requestBodies/%s", requestBodyName)
}

func responseBodyName(serviceName, methodName string) string {
	return fmt.Sprintf("%s%sResponse", serviceName, methodName)
}

func responseBodySchemaPath(responseBodyName string) string {
	return fmt.Sprintf("#/components/responses/%s", responseBodyName)
}
