package converter

import (
	"fmt"
	"strings"
)

func urlPath(microServiceName, protoServiceName, methodName string) string {
	return fmt.Sprintf("/%s/%s/%s", microServiceName, protoServiceName, methodName)
}

func payloadSchemaName(inputType string) string {
	payloadSchemaNameComponents := strings.Split(inputType, ".")
	return payloadSchemaNameComponents[len(payloadSchemaNameComponents)-1]
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
