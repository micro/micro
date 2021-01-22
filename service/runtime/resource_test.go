// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/runtime/resource_test.go

package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResources(t *testing.T) {

	// Namespace:
	assert.Equal(t, TypeNamespace, new(Namespace).Type())
	namespace, err := NewNamespace("")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, namespace)

	namespace, err = NewNamespace("test-namespace")
	assert.NoError(t, err)
	assert.NotNil(t, namespace)
	assert.Equal(t, TypeNamespace, namespace.Type())
	assert.Equal(t, "test-namespace", namespace.String())

	// NetworkPolicy:
	assert.Equal(t, TypeNetworkPolicy, new(NetworkPolicy).Type())
	networkPolicy, err := NewNetworkPolicy("", "", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, networkPolicy)

	networkPolicy, err = NewNetworkPolicy("test", "", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, networkPolicy)

	networkPolicy, err = NewNetworkPolicy("", "test", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, networkPolicy)

	networkPolicy, err = NewNetworkPolicy("ingress", "test", nil)
	assert.NoError(t, err)
	assert.NotNil(t, networkPolicy)
	assert.Equal(t, TypeNetworkPolicy, networkPolicy.Type())
	assert.Equal(t, "test.ingress", networkPolicy.String())
	assert.Len(t, networkPolicy.AllowedLabels, 1)

	networkPolicy, err = NewNetworkPolicy("ingress", "test", map[string]string{"foo": "bar", "bar": "foo"})
	assert.Len(t, networkPolicy.AllowedLabels, 2)

	// ResourceQuota:
	assert.Equal(t, TypeResourceQuota, new(ResourceQuota).Type())
	resourceQuota, err := NewResourceQuota("", "", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, resourceQuota)

	resourceQuota, err = NewResourceQuota("test", "", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, resourceQuota)

	resourceQuota, err = NewResourceQuota("", "test", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, resourceQuota)

	resourceQuota, err = NewResourceQuota("userquota", "test", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, networkPolicy)
	assert.Equal(t, TypeResourceQuota, resourceQuota.Type())
	assert.Equal(t, "test.userquota", resourceQuota.String())
	assert.Equal(t, 0, resourceQuota.Requests.CPU)
	assert.Equal(t, 2000, resourceQuota.Requests.Mem)
	assert.Equal(t, 5000, resourceQuota.Requests.Disk)
	assert.Equal(t, 0, resourceQuota.Limits.CPU)
	assert.Equal(t, 4000, resourceQuota.Limits.Mem)
	assert.Equal(t, 10000, resourceQuota.Limits.Disk)

	resourceQuota, err = NewResourceQuota("userquota", "test", &Resources{Mem: 777}, &Resources{Disk: 8888})
	assert.Equal(t, 0, resourceQuota.Requests.CPU)
	assert.Equal(t, 777, resourceQuota.Requests.Mem)
	assert.Equal(t, 0, resourceQuota.Requests.Disk)
	assert.Equal(t, 0, resourceQuota.Limits.CPU)
	assert.Equal(t, 0, resourceQuota.Limits.Mem)
	assert.Equal(t, 8888, resourceQuota.Limits.Disk)

	// Service:
	assert.Equal(t, TypeService, new(Service).Type())
	service, err := NewService("", "")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
	assert.Nil(t, service)

	service, err = NewService("test-service", "oldest")
	service.Metadata = map[string]string{"namespace": "testing"}
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, TypeService, service.Type())
	assert.Equal(t, "service://testing@test-service:oldest", service.String())
	assert.Equal(t, "oldest", service.Version)
}
