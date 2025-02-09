package namespace

import (
	"context"
	"errors"
	"strings"

	md "github.com/micro/micro/v5/service/context"
	"github.com/micro/micro/v5/service/registry"
	"github.com/micro/micro/v5/util/config"
)

const (
	// DefaultNamespace used by the server
	DefaultNamespace = "micro"
	// NamespaceKey is used to set/get the namespace from the context
	NamespaceKey = "Micro-Namespace"
)

const separator = ","

// FromContext gets the namespace from the context
func FromContext(ctx context.Context) string {
	// get the namespace which is set at ingress by micro web / api / proxy etc. The go-micro auth
	// wrapper will ensure the account making the request has the necessary issuer.
	ns, _ := md.Get(ctx, NamespaceKey)
	return ns
}

// ContextWithNamespace sets the namespace in the context
func ContextWithNamespace(ctx context.Context, ns string) context.Context {
	return md.Set(ctx, NamespaceKey, ns)
}

// List the namespaces for an environment
func List(env string) ([]string, error) {
	if len(env) == 0 {
		return nil, errors.New("Missing env value")
	}

	values, err := config.Get(config.Path("namespaces", env, "all"))
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return []string{registry.DefaultDomain}, nil
	}

	namespaces := strings.Split(values, separator)
	return append([]string{registry.DefaultDomain}, namespaces...), nil
}

// Add a namespace to an environment
func Add(namespace, env string) error {
	if len(env) == 0 {
		return errors.New("Missing env value")
	}
	if len(namespace) == 0 {
		return errors.New("Missing namespace value")
	}

	existing, err := List(env)
	if err != nil {
		return err
	}
	for _, ns := range existing {
		if ns == namespace {
			// the namespace already exists
			return nil
		}
	}

	values, _ := config.Get(config.Path("namespaces", env, "all"))
	if len(values) > 0 {
		values = strings.Join([]string{values, namespace}, separator)
	} else {
		values = namespace
	}

	return config.Set(config.Path("namespaces", env, "all"), values)
}

// Remove a namespace from an environment
func Remove(namespace, env string) error {
	if len(env) == 0 {
		return errors.New("Missing env value")
	}
	if len(namespace) == 0 {
		return errors.New("Missing namespace value")
	}
	if namespace == registry.DefaultDomain {
		return errors.New("Cannot remove the default namespace")
	}

	current, err := Get(env)
	if err != nil {
		return err
	}
	if current == namespace {
		err = Set(registry.DefaultDomain, env)
		if err != nil {
			return err
		}
	}

	existing, err := List(env)
	if err != nil {
		return err
	}

	var namespaces []string
	var found bool
	for _, ns := range existing {
		if ns == namespace {
			found = true
			continue
		}
		if ns == registry.DefaultDomain {
			continue
		}
		namespaces = append(namespaces, ns)
	}

	if !found {
		return errors.New("Namespace does not exists")
	}

	values := strings.Join(namespaces, separator)
	return config.Set(config.Path("namespaces", env, "all"), values)
}

// Set the current namespace for an environment
func Set(namespace, env string) error {
	if len(env) == 0 {
		return errors.New("Missing env value")
	}
	if len(namespace) == 0 {
		return errors.New("Missing namespace value")
	}

	existing, err := List(env)
	if err != nil {
		return err
	}

	var found bool
	for _, ns := range existing {
		if ns != namespace {
			continue
		}
		found = true
		break
	}

	if !found {
		return errors.New("Namespace does not exists")
	}

	return config.Set(config.Path("namespaces", env, "current"), namespace)
}

// Get the current namespace for an environment
func Get(env string) (string, error) {
	if len(env) == 0 {
		return "", errors.New("Missing env value")
	}

	if ns, err := config.Get(config.Path("namespaces", env, "current")); err != nil {
		return "", err
	} else if len(ns) > 0 {
		return ns, nil
	}

	return registry.DefaultDomain, nil
}
