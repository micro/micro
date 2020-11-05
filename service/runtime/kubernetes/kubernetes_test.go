// +build kubernetes

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
// Original source: github.com/micro/go-micro/v3/runtime/kubernetes/kubernetes_test.go
package kubernetes

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/micro/micro/v3/service/runtime"
	"github.com/stretchr/testify/assert"
)

func setupClient(t *testing.T) {
	files := []string{"token", "ca.crt"}
	for _, f := range files {
		cmd := exec.Command("kubectl", "get", "secrets", "-o",
			fmt.Sprintf(`jsonpath="{.items[?(@.metadata.annotations['kubernetes\.io/service-account\.name']=='micro-runtime')].data.%s}"`,
				strings.ReplaceAll(f, ".", "\\.")))
		if outp, err := cmd.Output(); err != nil {
			t.Fatalf("Failed to set k8s token %s", err)
		} else {
			outq := outp[1 : len(outp)-1]
			decoded, err := base64.StdEncoding.DecodeString(string(outq))
			if err != nil {
				t.Fatalf("Failed to set k8s token %s '%s'", err, outq)
			}
			if err := ioutil.WriteFile("/var/run/secrets/kubernetes.io/serviceaccount/"+f, decoded, 0755); err != nil {
				t.Fatalf("Error setting up k8s %s", err)
			}
		}

	}
	outp, err := exec.Command("kubectl", "config", "view", "-o", `jsonpath='{.clusters[?(@.name=="kind-kind")].cluster.server}'`).Output()
	if err != nil {
		t.Fatalf("Cannot find server for kind %s", err)
	}
	serverHost := string(outp)

	split := strings.Split(serverHost[9:len(serverHost)-1], ":")
	os.Setenv("KUBERNETES_SERVICE_HOST", split[0])
	os.Setenv("KUBERNETES_SERVICE_PORT", split[1])

}

func TestNamespaceCreateDelete(t *testing.T) {
	defer func() {
		exec.Command("kubectl", "-n", "foobar", "delete", "networkpolicy", "baz").Run()
		exec.Command("kubectl", "delete", "namespace", "foobar").Run()
	}()
	setupClient(t)
	r := NewRuntime()

	// Create a namespace
	testNamespace, err := runtime.NewNamespace("foobar")
	assert.NoError(t, err)
	if err := r.Create(testNamespace); err != nil {
		t.Fatalf("Unexpected error creating Namespace: %v", err)
	}

	// Check that the namespace exists
	if !namespaceExists(t, "foobar") {
		t.Fatalf("Namespace foobar not found")
	}

	// Create a networkpolicy:
	testNetworkPolicy, err := runtime.NewNetworkPolicy("baz", "foobar", nil)
	assert.NoError(t, err)
	if err := r.Create(testNetworkPolicy); err != nil {
		t.Fatalf("Unexpected error creating NetworkPolicy: %v", err)
	}

	// Check that the networkpolicy exists:
	if !networkPolicyExists(t, "foobar", "baz") {
		t.Fatalf("NetworkPolicy foobar.baz not found")
	}

	// Create a resourcequota:
	testResourceQuota, err := runtime.NewResourceQuota("caps", "foobar")
	assert.NoError(t, err)
	if err := r.Create(testResourceQuota); err != nil {
		t.Fatalf("Unexpected error creating ResourceQuota: %v", err)
	}

	// Check that the ResourceQuota exists:
	if !resourceQuotaExists(t, "foobar", "caps") {
		t.Fatalf("ResourceQuota foobar.caps not found")
	}

	// Tidy up
	if err := r.Delete(testResourceQuota); err != nil {
		t.Fatalf("Unexpected error deleting ResourceQuota: %v", err)
	}
	if resourceQuotaExists(t, "foobar", "caps") {
		t.Fatalf("ResourceQuota foobar.caps still exists")
	}
	if err := r.Delete(testNetworkPolicy); err != nil {
		t.Fatalf("Unexpected error deleting NetworkPolicy: %v", err)
	}
	if networkPolicyExists(t, "foobar", "baz") {
		t.Fatalf("NetworkPolicy foobar.baz still exists")
	}
	if err := r.Delete(testNamespace); err != nil {
		t.Fatalf("Unexpected error deleting Namespace: %v", err)
	}
	if namespaceExists(t, "foobar") {
		t.Fatalf("Namespace foobar still exists")
	}
}

func namespaceExists(t *testing.T, ns string) bool {
	cmd := exec.Command("kubectl", "get", "namespaces")
	outp, err := cmd.Output()
	if err != nil {
		t.Fatalf("Unexpected error listing namespaces %s", err)
	}
	exists, err := regexp.Match(ns+"\\s+Active", outp)
	if err != nil {
		t.Fatalf("Error listing namespaces %s", err)
	}
	return exists
}

func networkPolicyExists(t *testing.T, ns, np string) bool {
	cmd := exec.Command("kubectl", "-n", ns, "get", "networkpolicy")
	outp, err := cmd.Output()
	if err != nil {
		t.Fatalf("Unexpected error listing networkpolicies %s", err)
	}
	exists, err := regexp.Match(np, outp)
	if err != nil {
		t.Fatalf("Error listing networkpolicies %s", err)
	}
	return exists
}

func resourceQuotaExists(t *testing.T, ns, rq string) bool {
	cmd := exec.Command("kubectl", "-n", ns, "get", "resourcequota")
	outp, err := cmd.Output()
	if err != nil {
		t.Fatalf("Unexpected error listing resourcequotas %s", err)
	}
	exists, err := regexp.Match(rq, outp)
	if err != nil {
		t.Fatalf("Error listing resourcequotas %s", err)
	}
	return exists
}
