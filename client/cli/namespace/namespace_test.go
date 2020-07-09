package namespace

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/registry"
	"github.com/stretchr/testify/assert"
)

func TestNamespace(t *testing.T) {
	envName := fmt.Sprintf("test-%v", time.Now().UnixNano())
	namespace := "foo"

	t.Run("ListEmptyEnv", func(t *testing.T) {
		vals, err := List("")
		assert.Error(t, err, "Listing from a blank env should an error")
		assert.Nil(t, vals, "Listing from a blank env should not return namespaces")
	})

	t.Run("ListNewEnv", func(t *testing.T) {
		vals, err := List(envName)
		assert.Nilf(t, err, "Listing from a new env should not return an error")
		assert.Lenf(t, vals, 1, "Listing from a new env should return a single result")
		assert.Contains(t, vals, registry.DefaultDomain, "Listing from a new env should return the default namespace")
	})

	t.Run("AddEmptyEnv", func(t *testing.T) {
		err := Add("one", "")
		assert.Error(t, err, "Adding a namespace to an empty environment should return an error")
	})

	t.Run("AddEmptyNamespace", func(t *testing.T) {
		err := Add("", envName)
		assert.Error(t, err, "Adding an empty namespace to an environment should return an error")
	})

	t.Run("AddValidNamespace", func(t *testing.T) {
		err := Add(namespace, envName)
		assert.Nil(t, err, "Adding a valid namespace to an environment should not return an error")
	})

	t.Run("AddDuplicateNamespace", func(t *testing.T) {
		err := Add(namespace, envName)
		assert.Nil(t, err, "Adding a duplicate namespace to an environment should not return an error")
	})

	t.Run("ListPopulatedEnv", func(t *testing.T) {
		vals, err := List(envName)
		assert.Nilf(t, err, "Listing from a populated env should not return an error")
		assert.Lenf(t, vals, 2, "Listing from a populated env should return the correct number of results")
		assert.Contains(t, vals, registry.DefaultDomain, "Listing from a new env should return the default namespace")
		assert.Contains(t, vals, namespace, "Listing from a new env should return the added namespaces")
	})

	t.Run("GetBlankEnv", func(t *testing.T) {
		ns, err := Get("")
		assert.Error(t, err, "Getting from a blank env should an error")
		assert.Len(t, ns, 0, "Getting from a blank env should not return a namespace")
	})

	t.Run("GetUnsetEnv", func(t *testing.T) {
		ns, err := Get(envName)
		assert.Nil(t, err, "Getting an unset env should not error")
		assert.Equal(t, registry.DefaultDomain, ns, "Getting from an unset env should return the default namespace")
	})

	t.Run("SetBlankEnv", func(t *testing.T) {
		err := Set(namespace, "")
		assert.Error(t, err, "Setting the namespace of a blank env should an error")
	})

	t.Run("SetBlankNamespace", func(t *testing.T) {
		err := Set("", envName)
		assert.Error(t, err, "Setting a blank namespace should an error")
	})

	t.Run("SetInvalidNamespace", func(t *testing.T) {
		err := Set("notavalidns", envName)
		assert.Error(t, err, "Setting an unknown namespace should error")
	})

	t.Run("SetValidNamespace", func(t *testing.T) {
		err := Set(namespace, envName)
		assert.Nil(t, err, "Setting a valid namespace should not error")
	})

	t.Run("GetSetEnv", func(t *testing.T) {
		ns, err := Get(envName)
		assert.Nil(t, err, "Getting a set namespace should not error")
		assert.Equal(t, namespace, ns, "Getting a set namespace should return the correct value")
	})

	t.Run("RemoveEmptyEnv", func(t *testing.T) {
		err := Remove("one", "")
		assert.Error(t, err, "Removing a namespace from an empty environment should return an error")
	})

	t.Run("RemoveEmptyNamespace", func(t *testing.T) {
		err := Remove("", envName)
		assert.Error(t, err, "Removing an empty namespace from an environment should return an error")
	})

	t.Run("RemoveDefaultNamespace", func(t *testing.T) {
		err := Remove(registry.DefaultDomain, envName)
		assert.Error(t, err, "Removing the default namespace from an environment should return an error")
	})

	t.Run("RemoveSetNamespace", func(t *testing.T) {
		err := Remove(namespace, envName)
		assert.Error(t, err, "Removing a namespace which is currently set should error")
	})

	t.Run("SetDefaultNamespace", func(t *testing.T) {
		err := Set(registry.DefaultDomain, envName)
		assert.Nil(t, err, "Setting the default namespace should not error")
	})

	t.Run("RemoveValidNamespace", func(t *testing.T) {
		err := Remove(namespace, envName)
		assert.Nil(t, err, "Removing a valid namespace from an environment should not return an error")
	})

	t.Run("GetOverridenNamespace", func(t *testing.T) {
		ns, err := Get(envName)
		assert.Nil(t, err, "Getting an overriden namespace should not return an error")
		assert.Equal(t, registry.DefaultDomain, ns, "Getting an overriden namespace should return the correct value")
	})
}
