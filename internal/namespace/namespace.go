package namespace

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
)

const (
	// TODO: Move default namespace out of go-micro
	DefaultNamespace = auth.DefaultNamespace
	// RuntimeNamespace is the namespace which runtime services
	// such as the store and broker operate within. Any service
	// can read runtime services but writing is restricted.
	RuntimeNamespace = "runtime"
	// NamespaceKey is used to set/get the namespace from the
	// context
	NamespaceKey = "Micro-Namespace"
)

// NamespaceFromContext gets the namespace from the context
func NamespaceFromContext(ctx context.Context) string {
	if ns, ok := metadata.Get(ctx, NamespaceKey); ok {
		return ns
	}

	acc, err := auth.AccountFromContext(ctx)
	if err != nil || acc == nil {
		return DefaultNamespace
	}

	return acc.Namespace
}

var serviceTypes = []string{"api", "web", "service", "srv"}

// NamespaceFromService returns the namespace the service belongs to
func NamespaceFromService(name string) (string, error) {
	// joinKey is the key used to seperate the components of the service
	// name. '.' is the default, although '-' is also occasionally used.
	joinKey := "."
	if strings.ContainsAny(name, "-") {
		joinKey = "-"
	}

	// determine the type of service from the options in the serviceTypes
	// slice.
	var srvType string
	for _, t := range serviceTypes {
		// for when the srvType is in the middle of the name
		if strings.Contains(name, fmt.Sprintf("%v%v%v", joinKey, t, joinKey)) {
			srvType = t
			break
		}
		// for when the srvType is the first element in the name
		if strings.HasPrefix(name, fmt.Sprintf("%v%v", t, joinKey)) {
			srvType = t
			break
		}
	}

	// check to see if the service is a runtime service. This is true if the
	// namespace is the default namespace plus no serviceType was set.
	if len(srvType) == 0 && strings.HasPrefix(name, DefaultNamespace) {
		return RuntimeNamespace, nil
	} else if len(srvType) == 0 {
		return "", errors.BadRequest("go.micro.registry", "Missing service type in name")
	}

	// split the name into components and find the index of the srvType, since
	// all parts before this are the namespace and all parts after it are the
	// services alias.
	comps := strings.Split(name, joinKey)
	var typeIndex int
	for i, c := range comps {
		if c == srvType {
			typeIndex = i
			break
		}
	}

	// validate the typeIndex is not zero, causing an out of range error. This
	// would happen if no namespace is specified, in this case we use the default
	// one
	if typeIndex == 0 {
		return DefaultNamespace, nil
	}

	// the namespace is the components before the type, joined by the joinKey
	return strings.Join(comps[0:typeIndex], joinKey), nil
}
