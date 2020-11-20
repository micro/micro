package client

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/micro/micro/v3/service/runtime"

	"github.com/micro/micro/v3/internal/kubernetes/api"
	"github.com/micro/micro/v3/test/fakes"

	. "github.com/onsi/gomega"
)

func TestCreate(t *testing.T) {
	tcs := []struct {
		name         string
		namespace    string
		resource     *Resource
		expectedBody string
		expectedURL  string
	}{
		{
			name:      "deployment",
			namespace: "foo-bar-baz",
			resource: NewDeployment(&runtime.Service{
				Name:     "svc1",
				Version:  "latest",
				Source:   "source",
				Metadata: map[string]string{"foo": "bar", "hello": "world"},
			}, &runtime.CreateOptions{
				Command:   []string{"cmd", "arg"},
				Args:      []string{"arg1", "arg2"},
				Env:       []string{"FOO=BAR", "HELLO=WORLD"},
				Type:      "service",
				Image:     "DefaultImage",
				Namespace: DefaultNamespace,
				Resources: &runtime.Resources{
					CPU:  200,
					Mem:  200,
					Disk: 2000,
				},
				ServiceAccount: "serviceAcc",
			},
			),

			expectedBody: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "svc1-latest"
  namespace: "default"
  labels:
    micro: "service"
    name: "svc1"
    version: "latest"
  annotations:
    foo: "bar"
    hello: "world"
    name: "svc1"
    source: "source"
    version: "latest"
spec:
  replicas: 1
  selector:
    matchLabels:
      micro: "service"
      name: "svc1"
      version: "latest"
  template:
    metadata:
      labels:
        micro: "service"
        name: "svc1"
        version: "latest"
      annotations:
        foo: "bar"
        hello: "world"
        name: "svc1"
        source: "source"
        version: "latest"
    spec: 
      serviceAccountName: serviceAcc
      containers:
        - name: svc1
          env:
          - name: "FOO"
            value: "BAR"
          - name: "HELLO"
            value: "WORLD"
          args:
          - arg1
          - arg2
          command:
          - cmd
          - arg
          image: DefaultImage
          imagePullPolicy: IfNotPresent
          ports:
          - containerPort: 8080
            name: service-port
          readinessProbe:
            tcpSocket:
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 10
          resources:
            limits:
              memory: 200Mi
              cpu: 200m
              ephemeral-storage: 2000Mi
          volumeMounts: 
      volumes:
`,
			expectedURL: `example.com/apis/apps/v1/namespaces/foo-bar-baz/deployments/`,
		},
		{
			name:      "service",
			namespace: "foo-bar-baz",
			resource: NewService(&runtime.Service{
				Name:     "svc1",
				Version:  "latest",
				Source:   "source",
				Metadata: map[string]string{"foo": "bar", "hello": "world"},
			}, &runtime.CreateOptions{
				Command:   []string{"cmd", "arg"},
				Args:      []string{"arg1", "arg2"},
				Env:       []string{"FOO=BAR", "HELLO=WORLD"},
				Type:      "service",
				Image:     "DefaultImage",
				Namespace: DefaultNamespace,
				Resources: &runtime.Resources{
					CPU:  200,
					Mem:  200,
					Disk: 2000,
				},
				ServiceAccount: "serviceAcc",
			}),
			expectedBody: `
apiVersion: v1
kind: Service
metadata:
  name: "svc1"
  namespace: "default"
  labels:
    micro: "service"
    name: "svc1"
    version: "latest"
spec:
  selector:
    micro: "service"
    name: "svc1"
    version: "latest"
  ports:
  - name: "service-port"
    port: 8080
    protocol:
`,
			expectedURL: "example.com/api/v1/namespaces/foo-bar-baz/services/",
		},
		{
			name:      "secrets",
			namespace: "foo-bar-baz",
			resource: &Resource{
				Kind: "secret",
				Name: "svc1",
				Value: &Secret{
					Type: "Opaque",
					Data: map[string]string{
						"key1": "val1",
						"key2": "val2",
					},
					Metadata: &Metadata{
						Name:      "svc1",
						Namespace: "foo-bar-baz",
					},
				},
			},
			expectedBody: `
apiVersion: v1
kind: Secret
type: "Opaque"
metadata:
  name: "svc1"
  namespace: "foo-bar-baz"
  labels:
data:
  key1: "val1"
  key2: "val2"
`,
			expectedURL: "example.com/api/v1/namespaces/foo-bar-baz/secrets/",
		},
		{
			name:      "serviceaccount",
			namespace: "foo-bar-baz",
			resource: &Resource{
				Name: "svcacc",
				Kind: "serviceaccount",
				Value: &ServiceAccount{
					Metadata: &Metadata{
						Name:      "svcacc",
						Namespace: "foo-bar-baz",
					},
					ImagePullSecrets: []ImagePullSecret{
						{
							Name: "pullme",
						},
					},
				},
			},
			expectedBody: `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: "svcacc"
  labels:
imagePullSecrets:
- name: "pullme"
`,
			expectedURL: "example.com/api/v1/namespaces/foo-bar-baz/serviceaccounts/",
		},
		{
			name:      "networkpolicy",
			namespace: "foo-bar-baz",
			resource: &Resource{
				Kind: "networkpolicy",
				Value: NewNetworkPolicy("np1", "foo-bar-baz", map[string]string{
					"foo":   "bar",
					"hello": "world",
				})},
			expectedBody: `
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: "np1"
  namespace: "foo-bar-baz"
  labels:
spec:
  podSelector:
    matchLabels:
  ingress:
  - from: # Allow pods in this namespace to talk to each other
    - podSelector: {}
  - from: # Allow pods in the namespaces bearing the specified labels to talk to pods in this namespace:
    - namespaceSelector:
        matchLabels:
          foo: "bar"
          hello: "world"`,
			expectedURL: "example.com/apis/networking.k8s.io/v1/namespaces/foo-bar-baz/networkpolicies/",
		},
		{
			name:      "resourcequota",
			namespace: "foo-bar-baz",
			resource: &Resource{
				Kind: "resourcequota",
				Value: NewResourceQuota(&runtime.ResourceQuota{
					Name:      "rq1",
					Namespace: "foo-bar-baz",
					Requests: &runtime.Resources{
						CPU:  1000,
						Mem:  2000,
						Disk: 3000,
					},
					Limits: &runtime.Resources{
						CPU:  4000,
						Mem:  5000,
						Disk: 6000,
					},
				}),
			},
			expectedBody: `
apiVersion: v1
kind: ResourceQuota
metadata:
  name: "rq1"
  namespace: "foo-bar-baz"
  labels:
spec:
  hard:
    limits.memory: 5000Mi
    limits.cpu: 4000m
    limits.ephemeral-storage: 6000Mi
    requests.memory: 2000Mi
    requests.cpu: 1000m
    requests.ephemeral-storage: 3000Mi
`,
			expectedURL: "example.com/api/v1/namespaces/foo-bar-baz/resourcequotas/",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rt := fakes.FakeRoundTripper{}
			httpClient := &http.Client{
				Transport: &rt,
			}
			client := &client{
				opts: &api.Options{
					Host:      "example.com",
					Namespace: DefaultNamespace,
					Client:    httpClient,
				},
			}
			rt.RoundTripReturns(&http.Response{StatusCode: 200}, nil)

			g := NewWithT(t)
			err := client.Create(tc.resource, CreateNamespace(tc.namespace))
			g.Expect(err).ToNot(HaveOccurred())
			req := rt.RoundTripArgsForCall(0)
			g.Expect(req.URL.String()).To(Equal(tc.expectedURL))
			b, err := ioutil.ReadAll(req.Body)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(strings.TrimSpace(string(b))).To(Equal(strings.TrimSpace(tc.expectedBody)))
		})
	}

}
