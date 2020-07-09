package namespace

import (
	"errors"
	"strings"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/micro/v2/internal/config"
)

const seperator = ","

// List the namespaces for an environment
func List(env string) ([]string, error) {
	if len(env) == 0 {
		return nil, errors.New("Missing env value")
	}

	values, err := config.Get("namespaces", env, "all")
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return []string{registry.DefaultDomain}, nil
	}

	namespaces := strings.Split(values, seperator)
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

	values, _ := config.Get("namespaces", env, "all")
	if len(values) > 0 {
		values = strings.Join([]string{values, namespace}, seperator)
	} else {
		values = namespace
	}

	return config.Set(values, "namespaces", env, "all")
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
		return errors.New("Cannot remove the current namespace")
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

	values := strings.Join(namespaces, seperator)
	return config.Set(values, "namespaces", env, "all")
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

	return config.Set(namespace, "namespaces", env, "current")
}

// Get the current namespace for an environment
func Get(env string) (string, error) {
	if len(env) == 0 {
		return "", errors.New("Missing env value")
	}

	if ns, err := config.Get("namespaces", env, "current"); err != nil {
		return "", err
	} else if len(ns) > 0 {
		return ns, nil
	}

	return registry.DefaultDomain, nil
}
